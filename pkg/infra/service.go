package infra

import (
	"context"
	"getswizzle.io/swiz/pkg/common"
	"getswizzle.io/swiz/pkg/infra/aws"
	"getswizzle.io/swiz/pkg/infra/model"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Instancer interface {
	ListInstances() (map[string]model.TargetInstance, error)
}

type serviceDescriptor struct {
	description string
	instance    Instancer
}

type InfraService struct {
	services map[string]*serviceDescriptor
}

// NewInfraService creates a new infrastructure service
func NewInfraService() (*InfraService, error) {

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &InfraService{
		services: map[string]*serviceDescriptor{
			"AwsEc2": {
				description: "AWS EC2",
				instance:    aws.NewEc2(cfg),
			},
		},
	}, nil
}

// GetInstances returns all the instances from the specified service
func (s InfraService) GetInstances(serviceName string) (map[string]model.TargetInstance, error) {
	service := s.services[serviceName]
	if service == nil {
		return nil, common.NotFoundError{Subject: serviceName}
	}

	return service.instance.ListInstances()
}

// ListServices lists all of the services. A map by description and key will be returned
func (s InfraService) ListServices() map[string]string {
	serviceMap := map[string]string{}

	for k, v := range s.services {
		serviceMap[v.description] = k
	}

	return serviceMap
}
