package preprocessor

import (
	"fmt"
)

type node struct {
	name           string
	visited        bool
	visiting       bool
	dependentNames []string
	dependents     []*node
}

func (n *node) addDependent(dependent *node) {
	n.dependents = append(n.dependents, dependent)
}

func (n *node) hasCircularDependency() bool {
	n.visiting = true
	for _, dependent := range n.dependents {
		if dependent.visited {
			continue
		}
		if dependent.visiting || dependent.hasCircularDependency() {
			return true
		}
	}
	n.visiting = false
	n.visited = true
	return false
}

type DepTree struct {
	nodes map[string]*node
}

func NewDepTree() *DepTree {
	return &DepTree{
		nodes: map[string]*node{},
	}
}

func (d *DepTree) AddNodes(nodes map[string][]string) error {
	for name, dependencies := range nodes {
		if err := d.AddNode(name, dependencies); err != nil {
			return err
		}
	}
	return nil
}

func (d *DepTree) AddNode(name string, dependencies []string) error {

	if _, found := d.nodes[name]; found {
		return fmt.Errorf("node %s already exists", name)
	}

	d.nodes[name] = &node{
		name:           name,
		dependentNames: dependencies,
	}

	return nil
}

func (d *DepTree) Build() error {

	for _, n := range d.nodes {
		for _, dependName := range n.dependentNames {
			if dependNode, found := d.nodes[dependName]; found {
				n.addDependent(dependNode)
			} else {
				return fmt.Errorf("Could not find node for %s", dependName)
			}
		}
	}

	return nil
}

func (d *DepTree) GetCircularDependencies() []string {
	circularDependencies := []string{}

	for _, n := range d.nodes {
		if n.hasCircularDependency() {
			circularDependencies = append(circularDependencies, n.name)
		}
	}

	return circularDependencies
}

func (d *DepTree) Clear() {
	d.nodes = map[string]*node{}
}
