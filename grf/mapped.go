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

	panic("TODO")
}

func DeleteMapped(g *Graph, nt NodeType, key ...string) error {
	for _, k := range key {
		id, err := g.m.Map(nt, k)
		if err == nil {
			if err := Delete(g, id); err != nil {
				return err
			}
		} else if err != ErrNotFound {
			return err
		}
		if err := g.m.DeleteMapping(nt, k); err != nil {
			return err
		}
	}
	return nil
}
