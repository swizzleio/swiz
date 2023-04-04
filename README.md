# Swizzle

Swiz allows you to orchestrate multiple infrastructure-as-code (IaC) envionments. 

A common problem with distributed archectectures is there can be deploy time dependencies. For example, it's common to have a shared infrastructure stack. However if you deploy a stack the depends on the shared infrasturcture stack at the same time, you run the risk of the deployment failing due to a race condition between the stacks.

This application allows for the sequencing of these deployments.

## Roadmap

- [] Support for CloudFormation
- [] Support for the AWS CDK
- [] Support for terraform
- [] Azure and GCP support
- [] Cross cloud orchestration support
- [] Artisnal grill cheese sandwich delivery command