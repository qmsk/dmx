package main

import (
  log "github.com/Sirupsen/logrus"
  "github.com/qmsk/dmx/artnet"
  flags "github.com/jessevdk/go-flags"
  "fmt"
)

var options struct {
  Options

  Artnet  artnet.Config `group:"ArtNet"`
}

func main() {
  if args, err := flags.Parse(&options); err != nil {
    log.Fatalf("flags.Parse")
  } else if len(args) > 0 {
    log.Fatalf("Usage")
  } else {
    options.Setup()
  }

  var discoveryChan = make(chan artnet.Discovery)

  if artnetController, err := options.Artnet.Controller(); err != nil {
    log.Fatalf("artnet.Controller: %v", err)
  } else {
    log.Infof("artnet.Controller: %v", artnetController)

    artnetController.Start(discoveryChan)
  }

  dmx := artnet.Universe{
    0x00,
    0x00,
    0x00,
    0xff,
    0x00,
  }

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

        if err := node.SendDMX(outputPort.Address, dmx); err != nil {
          log.Errorf("Node %v: SendDMX: %v", node, err)
        } else {
          log.Infof("Node %v: SendDMX", node)
        }

        dmx[3]++
      }
    }
  }
}
