# Swizzle

Swiz allows you to orchestrate multiple infrastructure-as-code (IaC) environments. 

A common problem with distributed architectures is there can be deploy time dependencies. For example, it's common to 
have a shared infrastructure stack. However if you deploy a stack the depends on the shared infrastructure stack at the 
same time, you run the risk of the deployment failing due to a race condition between the stacks.

This application allows for the sequencing of these deployments.

## Mental Model

A stack represents a single service or deployment. It's the smallest possible unit. Typically a stack will have a single
git repo.

An environment consists of many stacks. Environments can be a product in production, development, test, PR, etc.

## Roadmap

The [product roadmap](https://github.com/orgs/swizzleio/projects/1) is kept here. Note, debt items are intermingled with
features.

## Development

### Prerequisites

#### Core skills
* [Golang](https://golang.org/doc/) is the language used.
* [AWS](https://docs.aws.amazon.com/index.html) is one of the cloud providers.
* [IaC concepts](https://en.wikipedia.org/wiki/Infrastructure_as_code) since that's core to what's being built.
* [Git](https://git-scm.com/book/en/v2) is used for source control.

#### Libraries used
* [AWS SDK for go](https://aws.github.io/aws-sdk-go-v2/docs/) for talking to AWS. Services used are:
  * [Cloudformation](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudformation)
  * [IAM](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/iam)
  * [Organizations](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/organizations)
  * [STS](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sts)
* [Testify](https://pkg.go.dev/github.com/stretchr/testify) for testing.
* [Urfave CLI V2](https://cli.urfave.org/v2/getting-started/) as the CLI library.
* [YAML](https://pkg.go.dev/gopkg.in/yaml.v3) for parsing YAML files.

## Building

The project uses Makefiles for building. Running `make` with no commands will print out the available commands. To build
the project, run `make build`. This will create a binary in the `out` directory.
