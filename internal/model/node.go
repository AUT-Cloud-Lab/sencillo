package model

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"slices"

	"k8s.io/apimachinery/pkg/api/resource"
)

var EdgeNodeList = []string{
	"custom-scheduler-debugging-m02",
	"custom-scheduler-debugging-m03",
	"custom-scheduler-debugging-m04",
}

type Node struct {
	id       types.UID
	name     string
	memory   *resource.Quantity
	cores    *resource.Quantity
	isOnEdge bool
}

func (node *Node) SetMemory(memory *resource.Quantity) {
	node.memory = memory
}

func (node *Node) SetCores(cores *resource.Quantity) {
	node.cores = cores
}

func NewNode(id types.UID, name string, memory *resource.Quantity, cpu *resource.Quantity) *Node {
	return &Node{
		id:       id,
		name:     name,
		memory:   memory,
		cores:    cpu,
		isOnEdge: checkIfOnEdge(name),
	}
}

func checkIfOnEdge(name string) bool {
	return slices.Contains(EdgeNodeList, name)
}

func (node *Node) ID() types.UID {
	return node.id
}

func (node *Node) Memory() *resource.Quantity {
	return node.memory
}

func (node *Node) Cores() *resource.Quantity {
	return node.cores
}
func (node *Node) Name() string {
	return node.name
}
func (node *Node) HasEnoughResourcesForPod(pod *Pod) bool {
	hasCpu := node.Cores().Cmp(*pod.Cores()) == 1
	hasMemory := node.Memory().Cmp(*pod.Memory()) == 1
	if !hasCpu {
		log.App.WithFields(logrus.Fields{
			"node_name":  node.Name(),
			"node_cores": node.Cores(),
			"is_on_edge": node.IsOnEdge(),
			"pod_cpu":    pod.Cores(),
		}).Info("is out of cpu")
	}

	if !hasMemory {
		log.App.WithFields(logrus.Fields{
			"node_name":   node.Name(),
			"node_memory": node.Memory(),
			"is_on_edge":  node.IsOnEdge(),
			"pod_memory":  pod.Memory(),
		}).Info("is out of memory")
	}

	return hasCpu && hasMemory
}

func (node *Node) IsOnEdge() bool {
	return node.isOnEdge
}
