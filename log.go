package main

import (
  log "github.com/Sirupsen/logrus"
)

type Options struct {
  Quiet     bool  `short:"q" long:"quiet"`
  Verbose   bool  `short:"v" long:"verbose"`
  Debug     bool  `long:"debug"`
}

func (options Options) Setup() {
  var logger = log.StandardLogger()

  if options.Debug {
    logger.Level = log.DebugLevel
  } else if options.Verbose {
    logger.Level = log.InfoLevel
  } else if options.Quiet {
    logger.Level = log.ErrorLevel
  } else {
    logger.Level = log.WarnLevel
  }
}
