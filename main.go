package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx/artnet"
	flags "github.com/jessevdk/go-flags"
	"time"
)

var options struct {
	Options

	Artnet artnet.Config `group:"ArtNet"`
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

	var artnetAddress = artnet.Address{}
	var dmxUniverse = artnet.Universe{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}

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

					if err := node.SendDMX(outputPort.Address, dmxUniverse); err != nil {
						log.Errorf("Node %v: SendDMX: %v", node, err)
					} else {
						log.Infof("Node %v: SendDMX", node)
					}
				}
			}
		}
	}()

	for range time.NewTicker(100 * time.Millisecond).C {
		dmxUniverse[3] += 10

		artnetController.SendDMX(artnetAddress, dmxUniverse)
	}
}
