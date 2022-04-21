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
	n, err := Mapped[T](g, nt, key)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		var zero T
		if n, err = Add(g, nt, zero); err != nil {
			return nil, err
		}
		if err := g.m.SetMapping(nt, key, n.ID); err != nil {
			return nil, err
		}
	}
	return update(g, n, fn)
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
