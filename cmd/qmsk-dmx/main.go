package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx/artnet"
	"github.com/SpComb/qmsk-dmx/heads"
	"github.com/SpComb/qmsk-web"
	flags "github.com/jessevdk/go-flags"
	colorful "github.com/lucasb-eyer/go-colorful"
)

var options struct {
	Options

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
		log.Infof("artnet.Discovery:")

		for _, node := range discovery.Nodes {
			fmt.Printf("%v:\n", node)

			config := node.Config()

			fmt.Printf("\tName: %v\n", config.Name)

			for i, inputPort := range config.InputPorts {
				fmt.Printf("\tInput %d: %v\n", i, inputPort.Address)
			}
			for i, outputPort := range config.OutputPorts {
				fmt.Printf("\tOutput %d: %v\n", i, outputPort.Address)

				// patch outputs
				universe := artnetController.Universe(outputPort.Address)
				// XXX: not safe
				hh.Output(heads.Universe(universe.Address().Integer()), universe)
			}
		}
	}

	log.Fatalf("artnet discovery ended")
}

func demo(hh *heads.Heads) {
	var intensity heads.Intensity = 1.0
	var hue float64 = 0.0

	for range time.NewTicker(100 * time.Millisecond).C {
		var color = colorful.Hsv(hue, 1.0, 1.0) // FastHappyColor()

		var headsColor = heads.ColorRGB{
			Red:   heads.Value(color.R),
			Green: heads.Value(color.G),
			Blue:  heads.Value(color.B),
		}

		hh.Each(func(head *heads.Head) {
			headParameters := head.Parameters()

			log.Debugf("head %v: intensity=%v color=%v", head, headParameters.Intensity, headParameters.Color)

			if headParameters.Color != nil {
				log.Debugf("head %v: Color %v @ %v", head, color, intensity)

				headParameters.Color.SetRGBIntensity(headsColor, intensity)

			} else if headParameters.Intensity != nil {
				log.Debugf("head %v: Intensity %v", head, intensity)

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
	if args, err := flags.Parse(&options); err != nil {
		log.Fatalf("flags.Parse")
	} else if len(args) > 0 {
		log.Fatalf("Usage")
	} else {
		options.Setup()
	}

	var artnetController *artnet.Controller

	if c, err := options.Artnet.Controller(); err != nil {
		log.Fatalf("artnet.Controller: %v", err)
	} else {
		log.Infof("artnet.Controller: %v", c)

		artnetController = c
	}

	// heads
	var headsHeads *heads.Heads

	if headsConfig, err := options.Heads.Config(options.Args.HeadsConfig); err != nil {
		log.Fatalf("heads.Config %v: %v", options.Args.HeadsConfig, err)
	} else if heads, err := options.Heads.Heads(headsConfig); err != nil {
		log.Fatalf("heads.Heads: %v", err)
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
		web.RoutePrefix("/api/", headsHeads.WebAPI()),
		options.Web.RouteStatic("/"),
	)
}
