package model

import (
	"fmt"
	"getswizzle.io/swiz/pkg/network"
)

type TargetInstance struct {
	Id        string
	Name      string
	Os        string
	Flavor    string
	Endpoints []network.Endpoint
}

func (t TargetInstance) String() string {
	return fmt.Sprintf("[%v] %v (%v)", t.Id, t.Name, t.Os)
}
