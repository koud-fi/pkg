package merge

type Getter interface{ Get(string) (any, bool) }

type GetterFunc func(string) (any, bool)

func (fn GetterFunc) Get(key string) (any, bool) { return fn(key) }
