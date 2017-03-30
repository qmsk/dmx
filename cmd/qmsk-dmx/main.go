package main

import (
	"time"

	"github.com/SpComb/qmsk-dmx/artnet"
	"github.com/SpComb/qmsk-dmx/heads"
	"github.com/SpComb/qmsk-dmx/logging"
	flags "github.com/jessevdk/go-flags"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/qmsk/go-web"
)

var options struct {
	Log logging.Option `long:"log"`

	Artnet artnet.Config `group:"ArtNet"`
	Heads  heads.Options `group:"Heads"`
	Web    web.Options   `group:"Web"`

	Demo bool `long:"demo" description:"Demo Effect"`

	Args struct {
		HeadsConfig string
	} `positional-args:"yes" required:"yes"`
}

// patch heads output universes on artnet discovery
func discovery(artnetController *artnet.Controller, hh *heads.Heads) {
	var discoveryChan = make(chan artnet.Discovery)

	artnetController.Start(discoveryChan)

	for discovery := range discoveryChan {
		for _, node := range discovery.Nodes {
			logging.Log.Infof("artnet.Discovery: %v:", node)

			config := node.Config()

			logging.Log.Infof("\tName: %v", config.Name)
			logging.Log.Infof("\tPorts: input=%d output=%d", len(config.InputPorts), len(config.OutputPorts))

			for i, inputPort := range config.InputPorts {
				logging.Log.Infof("\tInput %d: %v", i, inputPort.Address)
			}
			for i, outputPort := range config.OutputPorts {
				logging.Log.Infof("\tOutput %d: %v", i, outputPort.Address)

				// patch outputs
				universe := artnetController.Universe(outputPort.Address)

				var outputConfig = heads.OutputConfig{
					Universe: heads.Universe(outputPort.Address.Integer()),

					ArtNetNode: &config,
				}

				// XXX: not safe
				hh.Output(outputConfig, universe)
			}
		}
	}

	logging.Log.Fatalf("artnet discovery ended")
}

func demo(hh *heads.Heads) {
	var intensity heads.Intensity = 1.0
	var hue float64 = 0.0

	for range time.NewTicker(100 * time.Millisecond).C {
		var color = colorful.Hsv(hue, 1.0, 1.0) // FastHappyColor()

		var headsColor = heads.Color{
			Red:   heads.Value(color.R),
			Green: heads.Value(color.G),
			Blue:  heads.Value(color.B),
		}

		hh.Each(func(head *heads.Head) {
			headParameters := head.Parameters()

			logging.Log.Debugf("head %v: intensity=%v color=%v", head, headParameters.Intensity, headParameters.Color)

			if headParameters.Color != nil {
				logging.Log.Debugf("head %v: Color %v @ %v", head, color, intensity)

				headParameters.Color.SetIntensity(headsColor, intensity)

			} else if headParameters.Intensity != nil {
				logging.Log.Debugf("head %v: Intensity %v", head, intensity)

				headParameters.Intensity.Set(intensity)
			}
		})

		hh.Refresh()

		// animate
		intensity *= 0.95

		if intensity < 0.001 {
			intensity = 1.0
		}

		hue += 10.0

		if hue >= 360.0 {
			hue = 0.0
		}
	}
}

func main() {
	options.Log.Package = "main"

	if args, err := flags.Parse(&options); err != nil {
		logging.Log.Fatalf("flags.Parse")
	} else if len(args) > 0 {
		logging.Log.Fatalf("Usage")
	} else {
		logging.Setup(options.Log)
	}

	var artnetController *artnet.Controller

	if c, err := options.Artnet.Controller(); err != nil {
		logging.Log.Fatalf("artnet.Controller: %v", err)
	} else {
		logging.Log.Infof("artnet.Controller: %v", c)

		artnetController = c
	}

	// heads
	var headsHeads *heads.Heads

	if headsConfig, err := options.Heads.Config(options.Args.HeadsConfig); err != nil {
		logging.Log.Fatalf("heads.Config %v: %v", options.Args.HeadsConfig, err)
	} else if heads, err := options.Heads.Heads(headsConfig); err != nil {
		logging.Log.Fatalf("heads.Heads: %v", err)
	} else {
		headsHeads = heads
	}

	// artnet discovery to patch head outputs
	go discovery(artnetController, headsHeads)

	// animate heads
	if options.Demo {
		go demo(headsHeads)
	}

	// web
	options.Web.Server(
		options.Web.RouteEvents("/events", headsHeads.WebEvents()),
		options.Web.RouteAPI("/api/", headsHeads.WebAPI()),
		options.Web.RouteStatic("/"),
	)
}
