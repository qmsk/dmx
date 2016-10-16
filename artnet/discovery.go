package artnet

import (
  "net"
  "time"
)

type Discovery struct {
  Nodes       map[string]*Node
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

func (controller *Controller) discovery(pollChan chan pollEvent) {
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

      // TODO: timeout nodes

    case pollEvent := <-pollChan:
      if node := nodes[pollEvent.String()]; node != nil {
        node.discoveryTime = pollEvent.recvTime

        controller.log.Debugf("discovery refresh: %v", node)

      } else if node, err := controller.makeNode(pollEvent.srcAddr); err != nil {
        controller.log.Warnf("discovery %v: %v", pollEvent, err)

      } else {
        node.discoveryTime = pollEvent.recvTime

        controller.log.Debugf("discovery new: %v", node)

        nodes[pollEvent.String()] = node

        controller.update(nodes)
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

func (controller *Controller) update(nodes map[string]*Node) {
  var discovery = Discovery{
    Nodes:  make(map[string]*Node),
  }

  for name, node := range nodes {
    discovery.Nodes[name] = node
  }

  controller.discoveryState.Store(discovery)

  if controller.discoveryChan != nil {
    controller.discoveryChan <- discovery
  }
}

func (controller *Controller) Get() Discovery {
  if value := controller.discoveryState.Load(); value == nil {
    return Discovery{}
  } else {
    return value.(Discovery)
  }
}
