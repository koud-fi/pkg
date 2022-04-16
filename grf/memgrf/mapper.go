package memgrf

import "github.com/koud-fi/pkg/grf"

type mapper struct {
	locker
	Data map[grf.NodeType]map[string]grf.ID `json:"data"`
}

func NewMapper() grf.Mapper {
	return &mapper{Data: make(map[grf.NodeType]map[string]grf.ID)}
}

func (m *mapper) Map(nt grf.NodeType, key string) (grf.ID, error) {
	defer m.rlock()()
	if tm, ok := m.Data[nt]; ok {
		if id, ok := tm[key]; ok {
			return id, nil
		}
	}
	return 0, grf.ErrNotFound
}

func (m *mapper) SetMapping(nt grf.NodeType, key string, id grf.ID) error {
	defer m.lock()()
	if tm, ok := m.Data[nt]; ok {
		tm[key] = id
	} else {
		m.Data[nt] = map[string]grf.ID{key: id}
	}
	return nil
}

func (m *mapper) DeleteMapping(nt grf.NodeType, key ...string) error {
	defer m.lock()()
	if tm, ok := m.Data[nt]; ok {
		for _, key := range key {
			delete(tm, key)
		}
	}
	return nil
}
