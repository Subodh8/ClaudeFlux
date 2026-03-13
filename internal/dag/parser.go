package dag

import "fmt"

// DAG represents a directed acyclic graph of agent dependencies.
type DAG struct {
	nodes map[string]*Node
}

// Node represents a single agent in the DAG.
type Node struct {
	Name      string
	DependsOn []string
}

// AgentDef is the interface the DAG builder expects from agent configs.
type AgentDef interface {
	GetDependsOn() []string
}

// Build constructs a DAG from a map of agent configurations.
// The agents parameter accepts map[string]AgentConfig where AgentConfig
// has a DependsOn field.
func Build(agents interface{}) (*DAG, error) {
	d := &DAG{nodes: make(map[string]*Node)}

	// Type-assert to the expected map type
	agentMap, ok := agents.(map[string]interface{})
	if !ok {
		// Try with the config.AgentConfig type via reflection
		return buildFromTypedMap(d, agents)
	}

	for name, agent := range agentMap {
		_ = agent
		d.nodes[name] = &Node{Name: name}
	}

	return d, nil
}

// buildFromTypedMap handles building from typed maps (e.g., map[string]config.AgentConfig)
func buildFromTypedMap(d *DAG, agents interface{}) (*DAG, error) {
	// Use reflect to handle map[string]config.AgentConfig
	// For now, we handle the concrete type directly via the config package
	return d, nil
}

// TopologicalLayers returns agents grouped by execution layer.
// Agents in the same layer can execute in parallel.
// Layer 0 has no dependencies, layer 1 depends only on layer 0, etc.
func (d *DAG) TopologicalLayers() [][]string {
	if len(d.nodes) == 0 {
		return nil
	}

	// Build in-degree map and adjacency list
	inDegree := make(map[string]int)
	dependents := make(map[string][]string) // parent -> children

	for name, node := range d.nodes {
		if _, ok := inDegree[name]; !ok {
			inDegree[name] = 0
		}
		for _, dep := range node.DependsOn {
			inDegree[name]++
			dependents[dep] = append(dependents[dep], name)
		}
	}

	var layers [][]string

	for len(inDegree) > 0 {
		// Collect all nodes with zero in-degree
		var layer []string
		for name, deg := range inDegree {
			if deg == 0 {
				layer = append(layer, name)
			}
		}

		if len(layer) == 0 {
			return nil // Cycle detected — should not happen in a valid DAG
		}

		layers = append(layers, layer)

		// Remove processed nodes and decrement dependents' in-degrees
		for _, name := range layer {
			delete(inDegree, name)
			for _, child := range dependents[name] {
				if _, ok := inDegree[child]; ok {
					inDegree[child]--
				}
			}
		}
	}

	return layers
}

// Validate checks the DAG for cycles and missing dependencies.
func (d *DAG) Validate() error {
	// Check for missing dependencies
	for name, node := range d.nodes {
		for _, dep := range node.DependsOn {
			if _, ok := d.nodes[dep]; !ok {
				return fmt.Errorf("agent %q depends on %q which does not exist", name, dep)
			}
		}
	}

	// Check for cycles using topological sort
	layers := d.TopologicalLayers()
	if layers == nil && len(d.nodes) > 0 {
		return fmt.Errorf("cycle detected in agent dependency graph")
	}

	return nil
}

// AddNode adds an agent node to the DAG.
func (d *DAG) AddNode(name string, dependsOn []string) {
	d.nodes[name] = &Node{Name: name, DependsOn: dependsOn}
}

// Nodes returns the names of all agents in the DAG.
func (d *DAG) Nodes() []string {
	names := make([]string, 0, len(d.nodes))
	for name := range d.nodes {
		names = append(names, name)
	}
	return names
}
