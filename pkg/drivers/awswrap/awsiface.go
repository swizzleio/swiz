package awswrap

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

//go:generate mockery --name Iamer --filename iam_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type Iamer interface {
	iam.ListAccountAliasesAPIClient
}

//go:generate mockery --name Stser --filename sts_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type Stser interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

//go:generate mockery --name Orger --filename org_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type Orger interface {
	organizations.ListAccountsAPIClient
}

//go:generate mockery --name Cloudformationer --filename cloudformation_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type Cloudformationer interface {
	GetTemplateSummary(ctx context.Context, params *cloudformation.GetTemplateSummaryInput, optFns ...func(*cloudformation.Options)) (*cloudformation.GetTemplateSummaryOutput, error)
	CreateStack(ctx context.Context, params *cloudformation.CreateStackInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateStackOutput, error)
	DeleteStack(ctx context.Context, params *cloudformation.DeleteStackInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteStackOutput, error)
	CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error)
	DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error)
	ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error)
	DescribeStackResources(ctx context.Context, params *cloudformation.DescribeStackResourcesInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackResourcesOutput, error)

	cloudformation.DescribeChangeSetAPIClient
	cloudformation.DescribeStacksAPIClient
}

//go:generate mockery --name CfDescribeStacksPaginatorNewer --filename cloudformationpgnew_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type CfDescribeStacksPaginatorNewer func(client cloudformation.DescribeStacksAPIClient, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.DescribeStacksPaginatorOptions)) *cloudformation.DescribeStacksPaginator

//go:generate mockery --name CfDescribeStacksPaginator --filename cloudformationpg_mock.go --output ../../../mocks/ext/aws --outpkg mockaws
type CfDescribeStacksPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
}
