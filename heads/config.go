package heads

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
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

type ColorName string
type ColorChannel string

const (
	ColorChannelRed   = "red"
	ColorChannelGreen = "green"
	ColorChannelBlue  = "blue"
)

type ChannelType struct {
	Control   string       `json:",omitempty"`
	Intensity bool         `json:",omitempty"`
	Color     ColorChannel `json:",omitempty"`
}

func (channelType ChannelType) String() string {
	if channelType.Control != "" {
		return "control:" + channelType.Control
	}
	if channelType.Intensity {
		return "intensity"
	}
	if channelType.Color != "" {
		return "color:" + string(channelType.Color)
	}

	return ""
}

type HeadType struct {
	Vendor string
	Model  string
	Mode   string
	URL    string

	Channels []ChannelType
	Colors   map[ColorName]ColorRGB
}

func (headType HeadType) String() string {
	return fmt.Sprintf("%v/%v=%v", headType.Vendor, headType.Model, headType.Mode)
}

type HeadConfig struct {
	Type     string
	Universe Universe
	Address  dmx.Address
	Name     string
}

type Config struct {
	HeadTypes map[string]*HeadType

	Heads map[string]HeadConfig
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

func (options Options) Config(path string) (*Config, error) {
	var config = Config{
		HeadTypes: make(map[string]*HeadType),
	}

	if err := config.loadTypes(options.LibraryPath); err != nil {
		return nil, fmt.Errorf("loadTypes %v: %v", options.LibraryPath, err)
	}

	if _, err := load(&config, path); err != nil {
		return nil, err
	}

	return &config, nil
}
