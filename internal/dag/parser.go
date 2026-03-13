package dag

type DAG struct{}

func Build(agents interface{}) (*DAG, error) { return &DAG{}, nil }

func (d *DAG) TopologicalLayers() [][]string { return [][]string{} }
