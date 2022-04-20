package grf

/*
func (g *Graph) MappedNode(nt NodeType, key string, add bool) *Node {
	id, err := g.m.Map(nt, key)
	if err != nil {
		if add && err == ErrNotFound {
			n, err := g.AddNode(nt, nil)
			if err != nil {
				return &Node{err: err}
			}
			n.err = g.m.SetMapping(nt, key, n.id)
			return n
		}
		return &Node{err: err}
	}
	return g.Node(id)
}
*/
