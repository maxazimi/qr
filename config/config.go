package config

import (
	"encoding/gob"
	"fmt"
	"github.com/rapidloop/skv"
	"log"
	"sync"
)

const (
	key = "configs"
)

var (
	store    *skv.KVStore
	confInst *Config
	changed  bool
	mu       sync.Mutex
)

type Config struct {
	Map       map[string]any
	ThemeName string
}

func RegisterType(t any) {
	gob.Register(t)
}

func Set(ins *Config) {
	confInst = ins
	SetChanged()
}

func Get() *Config {
	return confInst
}

func SetChanged() {
	changed = true
}

func Load(dir string) (*Config, error) {
	mu.Lock()
	defer mu.Unlock()

	var err error
	dir += "/configs.db"

	store, err = skv.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open %s failed: %v", dir, err)
	}

	c := Config{}
	_ = store.Get(key, &c)

	if c.Map == nil {
		c.Map = make(map[string]any)
	}

	return &c, nil
}

func GetItem(key string) (any, error) {
	value, exists := confInst.Map[key]
	if !exists {
		return nil, fmt.Errorf("key[%s] does not exist", key)
	}
	return value, nil
}

func PutItem(key string, value any) {
	confInst.Map[key] = value
	SetChanged()
}

func Commit() error {
	mu.Lock()
	defer mu.Unlock()

	if err := store.Put(key, *confInst); err != nil {
		return err
	}

	log.Println("config saved")
	return nil
}

func Changed() bool {
	if changed {
		changed = false
		return true
	}
	return false
}

func (c *Config) Close() error {
	if err := Commit(); err != nil {
		return err
	}
	return store.Close()
}
