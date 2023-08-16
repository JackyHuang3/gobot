package global

import "sync"

type IStore interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte) error
}

var vstore *store

type store struct {
	cache sync.Map
}

func InitStore(vconf IConfig) (IStore, error) {
	vstore = &store{}
	return vstore, nil
}

func (p *store) Get(key string) ([]byte, bool) {
	val, ok := p.cache.Load(key)
	return val.([]byte), ok
}

func (p *store) Set(key string, value []byte) error {
	p.cache.Store(key, value)
	return nil
}
