# Swizzle

Swiz allows you to orchestrate multiple infrastructure-as-code (IaC) envionments. 

A common problem with distributed archectectures is there can be deploy time dependencies. For example, it's common to have a shared infrastructure stack. However if you deploy a stack the depends on the shared infrasturcture stack at the same time, you run the risk of the deployment failing due to a race condition between the stacks.

This application allows for the sequencing of these deployments.

## Mental Model

A stack represents a single service or deployment. It's the smallest possible unit. Typically a stack will have a single git repo.

An environment consists of many stacks. Environments can be a product in production, development, test, PR, etc.

## Roadmap

The [product roadmap](https://github.com/orgs/swizzleio/projects/1) is kept here. Note, debt items are intermingled with features.
