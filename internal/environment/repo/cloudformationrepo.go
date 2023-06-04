package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"strings"
	"time"
)

const CfPollTimeSec = 5

type CloudFormationRepo struct {
	client *cloudformation.Client
}

func NewCloudFormationRepo(config appconfig.AppConfig, enclave model.Enclave, provider *model.EncProvider) IacDeployer {
	cfg := provider.ToAwsConfig()

	return &CloudFormationRepo{
		client: cloudformation.NewFromConfig(cfg.GenerateConfig()),
	}
}

func (r *CloudFormationRepo) CreateStack(ctx context.Context, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {
	var stackInfo *model.StackInfo

	templateBody, templateUrl, err := r.templateOrUrl(template)
	if err != nil {
		return nil, fmt.Errorf("unable to get template body: %w", err)
	}

	templateResp, err := r.client.GetTemplateSummary(ctx, &cloudformation.GetTemplateSummaryInput{
		TemplateURL:  templateUrl,
		TemplateBody: templateBody,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get template summary: %w", err)
	}

	resources := templateResp.ResourceTypes

	state := model.StateDryRun
	reason := "Dry Run"
	details := ""
	if !dryRun {
		cfParams := r.generateParams(params, templateResp.Parameters)
		tags := r.generateTags(metadata)

		var resp *cloudformation.CreateStackOutput
		resp, err = r.client.CreateStack(ctx, &cloudformation.CreateStackInput{
			StackName:    &name,
			TemplateURL:  templateUrl,
			TemplateBody: templateBody,
			Parameters:   cfParams,
			Tags:         tags,
			OnFailure:    types.OnFailureDelete,
			Capabilities: []types.Capability{
				types.CapabilityCapabilityNamedIam,
			},
		})

		if resp != nil {
			details = *resp.StackId
		}
		reason = "Cloudformation CreateStack"
		state = model.StateCreating
	}

	if err != nil {
		return nil, fmt.Errorf("unable to create stack: %w", err)
	}

	stackInfo = &model.StackInfo{
		Name:       name,
		NextAction: model.NextActionCreate,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   state,
			Reason:  reason,
			Details: details,
		},
		Resources: resources,
	}

	return stackInfo, nil
}

func (r *CloudFormationRepo) DeleteStack(ctx context.Context, name string, dryRun bool) (*model.StackInfo, error) {

	var err error

	state := model.StateDryRun
	reason := "Dry Run"
	details := ""
	if !dryRun {
		_, err = r.client.DeleteStack(ctx, &cloudformation.DeleteStackInput{
			StackName: &name,
		})

		reason = "Cloudformation DeleteStack"
		state = model.StateDeleting
	}

	if err != nil {
		return nil, fmt.Errorf("unable to delete stack: %w", err)
	}

	stackInfo := &model.StackInfo{
		Name:       name,
		NextAction: model.NextActionDelete,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   state,
			Reason:  reason,
			Details: details,
		},
		Resources: []string{},
	}

	return stackInfo, nil
}

func (r *CloudFormationRepo) UpdateStack(ctx context.Context, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {

	// Get the current timestamp
	t := time.Now()
	timestamp := t.Format("20060102150405")
	changeSetName := fmt.Sprintf("Swz-%s-%s", name, timestamp)

	templateBody, templateUrl, err := r.templateOrUrl(template)
	if err != nil {
		return nil, fmt.Errorf("unable to get template body: %w", err)
	}

	templateResp, err := r.client.GetTemplateSummary(ctx, &cloudformation.GetTemplateSummaryInput{
		TemplateURL:  templateUrl,
		TemplateBody: templateBody,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get template summary: %w", err)
	}

	// Create change set
	cfParams := r.generateParams(params, templateResp.Parameters)
	tags := r.generateTags(metadata)
	_, err = r.client.CreateChangeSet(ctx, &cloudformation.CreateChangeSetInput{
		ChangeSetName: &changeSetName,
		StackName:     &name,
		TemplateURL:   templateUrl,
		TemplateBody:  templateBody,
		Parameters:    cfParams,
		Tags:          tags,
		Capabilities: []types.Capability{
			types.CapabilityCapabilityNamedIam,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create change set, %w", err)
	}

	// Wait for change set and fetch info
	notReady := true
	var resp *cloudformation.DescribeChangeSetOutput
	for notReady {
		resp, err = r.client.DescribeChangeSet(ctx, &cloudformation.DescribeChangeSetInput{
			ChangeSetName: &changeSetName,
			StackName:     &name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe change set, %w", err)
		}

		if len(resp.Changes) == 0 {
			// Empty change set, set dry run to true so we don't execute
			dryRun = true
			notReady = false
		}

		if resp.ExecutionStatus == types.ExecutionStatusAvailable {
			notReady = false
		} else {
			time.Sleep(CfPollTimeSec * time.Second)
		}
	}

	resourceChanges := []string{}
	for _, change := range resp.Changes {
		resourceChanges = append(resourceChanges, *change.ResourceChange.ResourceType)
	}

	state := model.StateDryRun
	reason := "Dry Run"
	details := ""

	// If this is a dry run, delete the change set
	if dryRun {
		_, err = r.client.DeleteChangeSet(ctx, &cloudformation.DeleteChangeSetInput{
			ChangeSetName: &changeSetName,
			StackName:     &name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to delete change set, %w", err)
		}
	} else {
		// Apply change set
		_, err = r.client.ExecuteChangeSet(ctx, &cloudformation.ExecuteChangeSetInput{
			ChangeSetName: &changeSetName,
			StackName:     &name,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to execute change set, %w", err)
		}

		state = model.StateUpdating
		reason = "Cloudformation UpdateStack"
		details = changeSetName
	}

	return &model.StackInfo{
		Name:       name,
		NextAction: model.NextActionUpdate,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   state,
			Reason:  reason,
			Details: details,
		},
		Resources: resourceChanges,
	}, nil
}

func (r *CloudFormationRepo) GetStackInfo(ctx context.Context, name string) (*model.StackInfo, error) {
	resp, err := r.client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: &name,
	})

	if err != nil {
		var apiError *smithy.GenericAPIError
		if errors.As(err, &apiError) && apiError.Code == "ValidationError" {
			return nil, apperr.NewNotFoundError("stack", name)
		}
		return nil, fmt.Errorf("fetching stack info: %w", err)
	}

	if len(resp.Stacks) > 0 {

		respResources, resErr := r.client.DescribeStackResources(ctx, &cloudformation.DescribeStackResourcesInput{
			StackName: &name,
		})
		if resErr != nil {
			return nil, fmt.Errorf("describe stack resources: %w", resErr)
		}

		resourceList := []string{}
		for _, res := range respResources.StackResources {
			resourceList = append(resourceList, *res.ResourceType)
		}

		return &model.StackInfo{
			Name:       name,
			NextAction: model.NextActionUpdate,
			DeployStatus: model.DeployStatus{
				Name:    name,
				State:   r.cfStatusToState(resp.Stacks[0].StackStatus),
				Reason:  r.strOrEmpty(resp.Stacks[0].StackStatusReason),
				Details: r.strOrEmpty(resp.Stacks[0].StackId),
			},
			Resources: resourceList,
		}, nil
	}

	return nil, fmt.Errorf("unable to find stack")
}

func (r *CloudFormationRepo) GetStackOutputs(ctx context.Context, name string) (map[string]string, error) {
	resp, err := r.client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: &name,
	})

	if err != nil {
		return nil, fmt.Errorf("unable to find resource: %w", err)
	}

	outputs := make(map[string]string)

	if len(resp.Stacks) > 0 {
		for _, item := range resp.Stacks[0].Outputs {
			outputs[*item.OutputKey] = r.strOrEmpty(item.OutputValue)
		}
	} else {
		return nil, fmt.Errorf("unable to find stack")
	}
	return outputs, nil
}

func (r *CloudFormationRepo) ListStacks(ctx context.Context, envName string) ([]model.StackInfo, error) {
	retVal := []model.StackInfo{}

	err := r.iterateAllStacks(ctx, func(stack types.Stack, tag types.Tag) (bool, error) {
		if *tag.Key == model.StackKeyEnvName && *tag.Value == envName {
			retVal = append(retVal, model.StackInfo{
				Name:       *stack.StackName,
				NextAction: model.NextActionUpdate,
				DeployStatus: model.DeployStatus{
					Name:    r.strOrEmpty(stack.StackName),
					State:   r.cfStatusToState(stack.StackStatus),
					Reason:  r.strOrEmpty(stack.StackStatusReason),
					Details: r.strOrEmpty(stack.StackId),
				},
				Resources: []string{}, // This is left blank because it's an expensive call
			})

			return true, nil
		}
		return false, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list stacks: %w", err)
	}

	return retVal, nil
}

func (r *CloudFormationRepo) ListEnvironments(ctx context.Context) ([]string, error) {
	// Store unique values
	uniqueValues := map[string]struct{}{}

	err := r.iterateAllStacks(ctx, func(stack types.Stack, tag types.Tag) (bool, error) {
		if *tag.Key == model.StackKeyEnvName {
			uniqueValues[*tag.Value] = struct{}{}

			return true, nil
		}
		return false, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	// Create slice of env names
	retVal := []string{}
	for value := range uniqueValues {
		retVal = append(retVal, value)
	}

	return retVal, nil
}

func (r *CloudFormationRepo) GetEnvironment(ctx context.Context, envName string) (*model.EnvironmentInfo, error) {

	stacks, err := r.ListStacks(ctx, envName)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	envState := model.StateComplete
	stackStatus := []string{}
	for _, stack := range stacks {
		if stack.DeployStatus.State != model.StateComplete {
			envState = envState.GetPriority(stack.DeployStatus.State)
			stackStatus = append(stackStatus, fmt.Sprintf("%v[%v]", stack.Name, stack.DeployStatus.State.String()))
		}
	}

	env := &model.EnvironmentInfo{
		EnvironmentName: envName,
		DeployStatus: model.DeployStatus{
			Name:    envName,
			State:   envState,
			Reason:  envState.String(),
			Details: strings.Join(stackStatus, ", "),
		},
		StackInfo: stacks,
	}

	return env, nil
}

func (r *CloudFormationRepo) IsEnvironmentInState(ctx context.Context, envName string, stacks []string, states []model.State) (bool, []string, error) {

	stackCompleteList := []string{}

	// Iterate over all stack names and fetch the state
	for _, stackName := range stacks {
		input := &cloudformation.DescribeStacksInput{StackName: &stackName}
		resp, err := r.client.DescribeStacks(ctx, input)
		if err != nil {
			return false, stackCompleteList, fmt.Errorf("failed to describe stack %s: %v", stackName, err)
		}

		// Print stack description
		for _, stack := range resp.Stacks {
			state := r.cfStatusToState(stack.StackStatus)
			for _, desiredState := range states {
				if state == desiredState {
					stackCompleteList = append(stackCompleteList, *stack.StackName)
				}
			}

			if state == model.StateFailed {
				return false, nil, fmt.Errorf("stack %v failed", stackName)
			}
		}
	}

	return len(stackCompleteList) == len(stacks), stackCompleteList, nil
}

func (r *CloudFormationRepo) templateOrUrl(template string) (templateBody *string, templateUrl *string, err error) {
	scheme, err := fileutil.GetScheme(template)
	if err != nil {
		return
	}

	if scheme == "file" {
		var b []byte
		b, err = fileutil.OpenUrl(template)
		if err != nil {
			return
		}

		str := string(b)
		templateBody = &str
	} else {
		templateUrl = &template
	}

	return
}

func (r *CloudFormationRepo) iterateAllStacks(ctx context.Context,
	stackFunc func(stack types.Stack, tag types.Tag) (bool, error)) error {
	describeStacksInput := &cloudformation.DescribeStacksInput{}
	paginator := cloudformation.NewDescribeStacksPaginator(r.client, describeStacksInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to get page when describing stacks: %w", err)
		}

		// Iterate over stacks and tags
		for _, stack := range page.Stacks {
			for _, tag := range stack.Tags {
				stop, ferr := stackFunc(stack, tag)
				if ferr != nil {
					return ferr
				}

				if stop {
					break
				}
			}
		}
	}
	return nil
}

func (r *CloudFormationRepo) generateTags(metadata map[string]string) []types.Tag {
	tags := []types.Tag{}
	for k, v := range metadata {
		tag := types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		tags = append(tags, tag)
	}
	return tags
}

func (r *CloudFormationRepo) strOrEmpty(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func (r *CloudFormationRepo) generateParams(params map[string]string, paramList []types.ParameterDeclaration) []types.Parameter {

	isParamNeeded := map[string]bool{}
	for _, param := range paramList {
		isParamNeeded[*param.ParameterKey] = true
	}

	cfParams := []types.Parameter{}
	for k, v := range params {
		if isParamNeeded[k] {
			cfParams = append(cfParams, types.Parameter{
				ParameterKey:   aws.String(k),
				ParameterValue: aws.String(v),
			})
		}
	}
	return cfParams
}

func (r *CloudFormationRepo) cfStatusToState(stackStatus types.StackStatus) model.State {
	switch stackStatus {
	case types.StackStatusCreateComplete:
		return model.StateComplete
	case types.StackStatusCreateInProgress:
		return model.StateCreating
	case types.StackStatusCreateFailed:
		return model.StateFailed
	case types.StackStatusDeleteComplete:
		return model.StateDeleted
	case types.StackStatusDeleteInProgress:
		return model.StateDeleting
	case types.StackStatusDeleteFailed:
		return model.StateFailed
	case types.StackStatusRollbackComplete:
		return model.StateFailed
	case types.StackStatusRollbackInProgress:
		return model.StateFailed
	case types.StackStatusRollbackFailed:
		return model.StateFailed
	case types.StackStatusUpdateComplete:
		return model.StateComplete
	case types.StackStatusUpdateInProgress:
		return model.StateUpdating
	case types.StackStatusUpdateRollbackComplete:
		return model.StateFailed
	case types.StackStatusUpdateRollbackInProgress:
		return model.StateFailed
	case types.StackStatusUpdateRollbackFailed:
		return model.StateFailed
	case types.StackStatusReviewInProgress:
		return model.StateUpdating
	case types.StackStatusUpdateCompleteCleanupInProgress:
		return model.StateUpdating
	case types.StackStatusUpdateFailed:
		return model.StateFailed
	case types.StackStatusUpdateRollbackCompleteCleanupInProgress:
		return model.StateFailed
	case types.StackStatusImportInProgress:
		return model.StateUpdating
	case types.StackStatusImportComplete:
		return model.StateComplete
	case types.StackStatusImportRollbackInProgress:
		return model.StateFailed
	case types.StackStatusImportRollbackFailed:
		return model.StateFailed
	case types.StackStatusImportRollbackComplete:
		return model.StateFailed
	default:
		return model.StateUnknown

	}
}
