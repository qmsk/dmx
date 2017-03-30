package heads

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/qmsk/dmx/logging"
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

func load(obj interface{}, path string) error {
	ext := filepath.Ext(path)

	switch ext {
	case ".toml":
		if err := loadToml(obj, path); err != nil {
			return fmt.Errorf("Load %T file %v: %v", obj, path, err)
		} else {
			return nil
		}
	default:
		return fmt.Errorf("Unknown %T file ext=%v: %v", obj, ext, path)
	}
}

type configMapper func(id []string) (configObject interface{}, err error)

// Load config from given path, using the given stat info
func loadsStat(path string, stat os.FileInfo, mapper configMapper, prefix []string, top bool) error {
	name := stat.Name()

	if stat.IsDir() {
		logging.Log.Debugf("heads:loadOne path=%v: dir prefix=%v", path, prefix)

		var id []string

		if top {
			id = prefix
		} else {
			id = append(id, prefix...)
			id = append(id, name)
		}

		return loadsDir(path, mapper, id)

	} else {
		// take basename *.ext part
		ext := filepath.Ext(name)
		name = name[:len(name)-len(ext)]

		var id []string

		if top {
			id = prefix
		} else {
			id = append(id, prefix...)
			id = append(id, name)
		}

		logging.Log.Debugf("heads:loadOne path=%v: file id=%v", path, id)

		if obj, err := mapper(id); err != nil {
			return err
		} else if err := load(obj, path); err != nil {
			return err
		} else {
			logging.Log.Infof("heads:loads path=%v: %T %v ", path, obj, id)

			return nil
		}
	}
}

// Recursively load multiple config files from a directory
func loadsDir(dirPath string, mapper configMapper, prefix []string) error {
	if files, err := ioutil.ReadDir(dirPath); err != nil {
		return fmt.Errorf("read dir %v: %v", dirPath, err)
	} else {
		for _, stat := range files {
			path := filepath.Join(dirPath, stat.Name())

			if !(stat.Mode().IsRegular() || stat.Mode().IsDir()) {
				logging.Log.Debugf("heads:loadsDir path=%v: skip irregular", path)

				continue
			}

			if strings.HasPrefix(stat.Name(), ".") {
				logging.Log.Debugf("heads:loadOne path=%v: skip dotfile", path)

				continue
			}

			if err := loadsStat(path, stat, mapper, prefix, false); err != nil {
				return err
			}
		}

		return nil
	}
}

// Load config from path, which may either be a file, or a directory to be loaded recursively
func loads(path string, mapper configMapper) error {
	if stat, err := os.Stat(path); err != nil {
		return err
	} else {
		return loadsStat(path, stat, mapper, nil, true)
	}
}

type Config struct {
	HeadTypes map[TypeID]*HeadType
	Colors    map[ColorID]*Color

	Heads  map[HeadID]*HeadConfig
	Groups map[GroupID]*GroupConfig

	Presets map[PresetID]*PresetConfig
}

func (config *Config) loadTypes(path string) error {
	return loads(path, func(id []string) (interface{}, error) {
		var typeID = TypeID(filepath.Join(id...))
		var headType = new(HeadType)

		config.HeadTypes[typeID] = headType

		return headType, nil
	})
}

func (config *Config) load(path string) error {
	return loads(path, func(id []string) (interface{}, error) {
		if len(id) == 0 {
			return config, nil
		}
		switch id[0] {
		case "colors":
			if len(id) == 1 {
				return &config.Colors, nil
			} else {
				var color = new(Color)

				config.Colors[ColorID(filepath.Join(id[1:]...))] = color

				return color, nil
			}

		case "heads":
			if len(id) == 1 {
				return &config.Heads, nil
			} else {
				var head = new(HeadConfig)

				config.Heads[HeadID(filepath.Join(id[1:]...))] = head

				return head, nil
			}

		case "groups":
			if len(id) == 1 {
				return &config.Groups, nil
			} else {
				var group = new(GroupConfig)

				config.Groups[GroupID(filepath.Join(id[1:]...))] = group

				return group, nil
			}

		case "presets":
			if len(id) == 1 {
				return &config.Presets, nil
			} else {
				var preset = new(PresetConfig)

				config.Presets[PresetID(filepath.Join(id[1:]...))] = preset

				return preset, nil
			}

		default:
			return nil, fmt.Errorf("Bad config path: %v", id)
		}
	})

}

// map relative Head.Type= references
func (config *Config) mapTypes() error {
	// clone to ColorMap without pointers
	var colors = make(ColorMap)
	for colorID, color := range config.Colors {
		colors[colorID] = *color
	}

	// inherit colors
	for _, headType := range config.HeadTypes {
		if !headType.IsColor() {
			continue
		}

		if headType.Colors == nil {
			headType.Colors = make(ColorMap)
		}

		// each headType has its own copy
		headType.Colors.Merge(colors)
	}

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
		HeadTypes: make(map[TypeID]*HeadType),
		Colors:    make(map[ColorID]*Color),
		Heads:     make(map[HeadID]*HeadConfig),
		Groups:    make(map[GroupID]*GroupConfig),
		Presets:   make(map[PresetID]*PresetConfig),
	}

	for _, libraryPath := range options.LibraryPath {
		if err := config.loadTypes(libraryPath); err != nil {
			return nil, fmt.Errorf("loadTypes %v: %v", libraryPath, err)
		}
	}

	if err := config.load(path); err != nil {
		return nil, err
	}

	if err := config.mapTypes(); err != nil {
		return nil, err
	}

	return &config, nil
}
