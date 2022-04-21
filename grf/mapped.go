package grf

func Mapped[T any](g *Graph, nt NodeType, key string) (*Node[T], error) {

	// ???

	panic("TODO")
}

func SetMapped[T any](
	g *Graph, nt NodeType, key string, fn func(T) (T, error),
) (*Node[T], error) {

	// ???

	panic("TODO")
}

func DeleteMapped(g *Graph, nt NodeType, key ...string) error {

	// ???

	panic("TODO")
}

/*
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
*/
