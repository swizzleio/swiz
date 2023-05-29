package awswrap

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type CfWrap struct {
	client *cloudformation.Client
}

func NewCfWrapper(cfg AwsConfig) *CfWrap {
	return &CfWrap{
		client: cloudformation.NewFromConfig(cfg.GenerateConfig()),
	}
}

func (c *CfWrap) GetOutput(stackName string, resourceName string) (string, error) {
	result, err := c.client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{
		StackName: &stackName})

	if err != nil {
		return "", fmt.Errorf("unable to find resource: %w", err)
	}

	if len(result.Stacks) > 0 {
		for _, item := range result.Stacks[0].Outputs {
			//fmt.Println(*item.OutputKey, " : ", *item.OutputValue, " - ", resourceName)
			if *item.OutputKey == resourceName {
				return *item.OutputValue, nil
			}
		}
	} else {
		return "", fmt.Errorf("unable to find stack")
	}

	return "", nil
}

func (c *CfWrap) GetOutputs(stackName string) (map[string]string, error) {
	result, err := c.client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{
		StackName: &stackName})

	if err != nil {
		return nil, fmt.Errorf("unable to find resource: %w", err)
	}

	outputs := make(map[string]string)

	if len(result.Stacks) > 0 {
		for _, item := range result.Stacks[0].Outputs {
			outputs[*item.OutputKey] = *item.OutputValue
		}
	} else {
		return nil, fmt.Errorf("unable to find stack")
	}

	return outputs, nil
}

func (c *CfWrap) ListStacks() ([]string, error) {
	result, err := c.client.ListStacks(context.TODO(), &cloudformation.ListStacksInput{})

	if err != nil {
		return nil, fmt.Errorf("unable to list stacks: %w", err)
	}

	stacks := make([]string, 0)

	for _, stack := range result.StackSummaries {
		stacks = append(stacks, *stack.StackName)
	}

	return stacks, nil
}

func (c *CfWrap) CreateStack(stackName string, templateBody string, parameters []types.Parameter) error {
	_, err := c.client.CreateStack(context.TODO(), &cloudformation.CreateStackInput{
		StackName:    &stackName,
		TemplateBody: &templateBody,
		Parameters:   parameters,
	})

	if err != nil {
		return fmt.Errorf("unable to create stack: %w", err)
	}

	return nil
}

func (c *CfWrap) DeleteStack(stackName string) error {
	_, err := c.client.DeleteStack(context.TODO(), &cloudformation.DeleteStackInput{
		StackName: &stackName,
	})

	if err != nil {
		return fmt.Errorf("unable to delete stack: %w", err)
	}

	return nil
}

func (c *CfWrap) WaitForStack(stackName string, timeoutMin float64) error {
	waiter := cloudformation.NewStackCreateCompleteWaiter(c.client)
	err := waiter.Wait(context.TODO(), &cloudformation.DescribeStacksInput{StackName: &stackName},
		time.Duration(timeoutMin)*time.Minute)

	if err != nil {
		return fmt.Errorf("unable to wait for stack: %w", err)
	}

	return nil
}
