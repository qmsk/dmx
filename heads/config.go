package heads

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

func loadToml(obj interface{}, path string) error {
	if meta, err := toml.DecodeFile(path, obj); err != nil {
		return err
	} else if len(meta.Undecoded()) > 0 {
		return fmt.Errorf("Invalid keys in toml: %#v", meta.Undecoded())
	} else {
		return nil
	}
}

func load(obj interface{}, path string) (string, error) {
	base := filepath.Base(path)
	ext := filepath.Ext(path)

	name := base[:len(base)-len(ext)]

	switch ext {
	case ".toml":
		if err := loadToml(obj, path); err != nil {
			return name, err
		} else {
			return name, nil
		}
	default:
		return name, fmt.Errorf("Unknown %T file ext=%v: %v", obj, ext, path)
	}
}

type Config struct {
	HeadTypes map[string]*HeadType

	Heads  map[HeadID]*HeadConfig
	Groups map[GroupID]GroupConfig
}

func (config *Config) loadTypes(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var headType HeadType

		log.Debugf("heads:Config.loadTypes %v: %v mode=%v", rootPath, path, info.Mode())

		if !info.Mode().IsRegular() {
			return nil
		}

		relPath := path[len(rootPath):]
		if relPath[0] == '/' {
			relPath = relPath[1:]
		}
		dir, name := filepath.Split(relPath)

		if basename, err := load(&headType, path); err != nil {
			return err
		} else {
			name = filepath.Join(dir, basename)
		}

		log.Infof("heads:Config.loadTypes %v: HeadType %v", path, name)

		config.HeadTypes[name] = &headType

		return nil
	})
}

// map relative Head.Type= references
func (config *Config) mapTypes() error {
	for headID, headConfig := range config.Heads {
		if headType, exists := config.HeadTypes[headConfig.Type]; !exists {
			return fmt.Errorf("heads.%s: Invalid Head.Type=%v", headID, headConfig.Type)
		} else {
			headConfig.headType = headType
		}
	}

	return nil
}

func (options Options) Config(path string) (*Config, error) {
	var config = Config{
		HeadTypes: make(map[string]*HeadType),
		Groups:    make(map[GroupID]GroupConfig),
	}

	if err := config.loadTypes(options.LibraryPath); err != nil {
		return nil, fmt.Errorf("loadTypes %v: %v", options.LibraryPath, err)
	}

	if _, err := load(&config, path); err != nil {
		return nil, err
	}

	if err := config.mapTypes(); err != nil {
		return nil, err
	}

	return &config, nil
}
