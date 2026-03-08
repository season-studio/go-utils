package config

import (
	"fmt"
	"os"
	"season-studio/go-utils/log"
	"sync"

	"gopkg.in/yaml.v3"
)

type _RegistryEntry struct {
	key    string
	initFn func(cfg []byte) any
	data   any
}

var (
	locker        sync.Mutex
	registryList  []*_RegistryEntry
	loadedActions []func()
)

func RegisterEntry[T any](key string, fnInit func(v *T)) {
	stubFn := func(cfg []byte) (retVal any) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Exception raised when loading configuration of \"%s\"\n\t%v", key, err)
				retVal = nil
			}
		}()
		var data T
		err := yaml.Unmarshal(cfg, &data)
		if err != nil {
			panic(err)
		}
		fnInit(&data)
		return &data
	}

	locker.Lock()
	defer locker.Unlock()

	registryList = append(registryList, &_RegistryEntry{
		key:    key,
		initFn: stubFn,
		data:   nil,
	})
}

func OnLoaded(fn func()) {
	locker.Lock()
	defer locker.Unlock()

	loadedActions = append(loadedActions, fn)
}

func LoadFromBytes(bytes []byte) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("%v", err)
			log.Errorf("Exception raised when loading configuration\n\t%v", err)
			log.Flush()
		}
	}()

	locker.Lock()
	defer locker.Unlock()

	var tempMap map[string]any
	err := yaml.Unmarshal(bytes, &tempMap)
	if err != nil {
		panic(err)
	}

	for _, entry := range registryList {
		if len(entry.key) == 0 {
			data := entry.initFn(bytes)
			entry.data = data
		} else {
			var cfgBytes []byte
			var err error
			if cfg, ok := tempMap[entry.key]; ok {
				cfgBytes, err = yaml.Marshal(cfg)
				if err != nil {
					log.Errorf("Exception raised when marshal configuration of \"%s\"\n\t%v", entry.key, err)
					continue
				}
			} else {
				cfgBytes = []byte("")
			}
			data := entry.initFn(cfgBytes)
			entry.data = data
		}
	}

	for _, fn := range loadedActions {
		fn()
	}

	return nil
}

func Load(filePath string) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Exception raised when loading configuration from file %s\n\t%v", filePath, err)
		}
	}()

	filedata, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	LoadFromBytes(filedata)
}

func Save() (retVal []byte, retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("%v", err)
			retVal = nil
			log.Errorf("Exception raised when saving configuration\n\t%v", err)
		}
	}()

	locker.Lock()
	defer locker.Unlock()

	storage := make(map[string]any)
	for _, entry := range registryList {
		if entry.data == nil {
			continue
		}
		if len(entry.key) == 0 {
			bytes, err := yaml.Marshal(entry.data)
			if err != nil {
				log.Errorf("Exception raised when marshal configuration of \"%s\" for saving\n\t%v", entry.key, err)
				continue
			}
			temp := make(map[string]any)
			err = yaml.Unmarshal(bytes, &temp)
			if err != nil {
				log.Errorf("Exception raised when unmarshal configuration of \"%s\" for saving\n\t%v", entry.key, err)
				continue
			}
			for key, data := range temp {
				storage[key] = data
			}
		} else {
			storage[entry.key] = entry.data
		}
	}

	bytes, err := yaml.Marshal(storage)
	if err != nil {
		panic(err)
	}

	return bytes, nil
}

func Get[T any]() *T {
	locker.Lock()
	defer locker.Unlock()

	for _, entry := range registryList {
		data, ok := entry.data.(*T)
		if ok {
			return data
		}
	}
	return nil
}

func GetByKey(key string) any {
	locker.Lock()
	defer locker.Unlock()

	for _, entry := range registryList {
		if entry.key == key {
			return entry.data
		}
	}
	return nil
}
