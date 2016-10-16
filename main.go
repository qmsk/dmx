package main

import (
  log "github.com/Sirupsen/logrus"
  "github.com/qmsk/dmx/artnet"
  flags "github.com/jessevdk/go-flags"
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

  for discovery := range discoveryChan {
    log.Infof("artnet.Discovery:")

    for _, node := range discovery.Nodes {
      log.Infof("\t%v", node)
    }
  }
}
