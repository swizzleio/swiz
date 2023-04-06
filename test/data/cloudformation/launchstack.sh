#!/bin/bash

# Create stacks
aws cloudformation create-stack --template-body file://bootstrap.yaml --stack-name swiz-boot --capabilities CAPABILITY_NAMED_IAM
aws cloudformation wait stack-create-complete --stack-name swiz-boot
LAMBDA_ARN=`aws cloudformation describe-stacks --stack-name swiz-boot --query "Stacks[0].Outputs[?OutputKey=='SleepTestFunctionArn'].OutputValue" --output text`
aws cloudformation create-stack --template-body file://sleepstack.yaml --stack-name swiz-sleep --capabilities CAPABILITY_NAMED_IAM --parameters ParameterKey=SleepTestFunctionArn,ParameterValue=$LAMBDA_ARN

# Delete stacks
aws cloudformation wait stack-create-complete --stack-name swiz-sleep
aws cloudformation delete-stack --stack-name swiz-boot
aws cloudformation delete-stack --stack-name swiz-sleep