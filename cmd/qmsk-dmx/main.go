package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx/artnet"
	"github.com/SpComb/qmsk-dmx/heads"
	flags "github.com/jessevdk/go-flags"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/qmsk/e2/web"
)

var options struct {
	Options

	Artnet artnet.Config `group:"ArtNet"`
	Heads  heads.Options `group:"Heads"`
	Web    web.Options   `group:"Web"`

	Args struct {
		HeadsConfig string
	} `positional-args:"yes" required:"yes"`
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
	var discoveryChan = make(chan artnet.Discovery)

	if c, err := options.Artnet.Controller(); err != nil {
		log.Fatalf("artnet.Controller: %v", err)
	} else {
		log.Infof("artnet.Controller: %v", c)

		c.Start(discoveryChan)

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

	// patch heads output universes on artnet discovery
	go func() {
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
				}
			}

			// patch outputs
			for address, universe := range artnetController.Universes() {
				// XXX: not safe
				headsHeads.Output(heads.Universe(address.Integer()), universe)
			}
		}
	}()

	// web
	go options.Web.Server(
		web.RoutePrefix("/api/", headsHeads.WebAPI()),
	)

	// animate heads
	var intensity heads.Intensity = 1.0
	var hue float64 = 0.0

	for range time.NewTicker(100 * time.Millisecond).C {
		var color = colorful.Hsv(hue, 1.0, 1.0) // FastHappyColor()

		var headsColor = heads.ColorRGB{
			R: heads.Value(color.R),
			G: heads.Value(color.G),
			B: heads.Value(color.B),
		}

		headsHeads.Each(func(head *heads.Head) {
			headIntensity := head.Intensity()
			headColor := head.Color()

			log.Debugf("head %v: intensity=%v color=%v", head, headIntensity.Get(), headColor.Exists())

			if headColor.Exists() {
				log.Debugf("head %v: Color %v @ %v", head, color, intensity)

				headColor.SetRGBIntensity(headsColor, intensity)

			} else if headIntensity.Exists() {
				log.Debugf("head %v: Intensity %v", head, intensity)

				headIntensity.Set(intensity)
			}
		})

		headsHeads.Refresh()

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
