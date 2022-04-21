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
	if n, err := Mapped[T](g, nt, key); err == nil {
		return update(g, n, fn)
	} else if err != ErrNotFound {
		return nil, err
	}

	// ???

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
