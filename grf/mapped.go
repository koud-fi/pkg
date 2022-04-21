package grf

func Mapped[T any](g *Graph, nt NodeType, key string) (*Node[T], error) {
	id, err := g.m.Map(nt, key)
	if err != nil {
		return nil, err
	}
	return Lookup[T](g, id)
}

func SetMapped[T any](
	g *Graph, nt NodeType, key string, fn func(T) (T, error),
) (*Node[T], error) {

	/*
		n, err := g.AddNode(nt, nil)
		if err != nil {
			return &Node{err: err}
		}
		n.err = g.m.SetMapping(nt, key, n.id)
		return n
	*/

	panic("TODO")
}

func DeleteMapped(g *Graph, nt NodeType, key ...string) error {

	// TODO: delete the actual nodes

	return g.m.DeleteMapping(nt, key...)
}
