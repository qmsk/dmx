package artnet

import (
  "net"
  "time"
)

type discoveryEvent struct {
  recvTime    time.Time
  srcAddr     *net.UDPAddr
  pollReply   ArtPollReply
}

func (event discoveryEvent) String() string {
  return event.srcAddr.String()
}

func (controller *Controller) sendPoll(addr *net.UDPAddr) error {
  return controller.transport.send(addr, ArtPoll{
    ArtHeader: ArtHeader{
      ID: ARTNET,
      OpCode: OpPoll,
    },
    ProtVer: ProtVer,
  })
}

func (controller *Controller) discovery(discoveryChan chan discoveryEvent) {
  var ticker = time.NewTicker(controller.discoveryInterval)
  var nodes = make(map[string]*Node)

  if err := controller.sendPoll(controller.discoveryAddr); err != nil {
    controller.log.Fatalf("discovery: sendPoll: %v", err)
  }

  for {
    select {
    case <-ticker.C:
      controller.log.Debug("discovery: tick...")

      if err := controller.sendPoll(controller.discoveryAddr); err != nil {
        controller.log.Fatalf("discovery: sendPoll: %v", err)
      }

    case event := <-discoveryChan:
      if node := nodes[event.String()]; node != nil {
        node.discoveryTime = event.recvTime

        controller.log.Debugf("discovery: %v", node)

      } else if node, err := controller.makeNode(event.srcAddr); err != nil {
        controller.log.Warnf("discovery %v: %v", event, err)

      } else {
        node.discoveryTime = event.recvTime

        controller.log.Infof("discovery: %v", event, node)

        nodes[event.String()] = node
      }
    }
  }
}

func (controller *Controller) makeNode(addr *net.UDPAddr) (*Node, error) {
  var node = Node{
    log:  controller.log.WithField("node", addr.String()),

    transport:  controller.transport,
    addr:       addr,
  }

  return &node, nil
}
