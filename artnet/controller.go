package artnet

import (
  "bytes"
  "fmt"
  "net"
  "encoding/binary"
  log "github.com/Sirupsen/logrus"
)

type Config struct {
  Listen    string    `long:"artnet-listen" value-name:"ADDR" default:"0.0.0.0"`
  Broadcast string    `long:"artnet-broadcast" value-name:"ADDR" default:"255.255.255.255"`
}

func (config Config) Controller() (*Controller, error) {
  var controller = Controller{
    log:  log.WithFields(log.Fields{"prefix": "artnet:Controller"}),
  }

  listenAddr  := net.JoinHostPort(config.Listen, fmt.Sprintf("%d", Port))
  broadcastAddr := net.JoinHostPort(config.Broadcast, fmt.Sprintf("%d", Port))

  if udpAddr, err := net.ResolveUDPAddr("udp", listenAddr); err != nil {
    return nil, err
  } else if udpConn, err := net.ListenUDP("udp", udpAddr); err != nil {
    return nil, err
  } else {
    controller.udpConn = udpConn
  }

  if udpAddr, err := net.ResolveUDPAddr("udp", broadcastAddr); err != nil {
    return nil, err
  } else {
    controller.udpAddr = udpAddr
  }

  return &controller, nil
}

type Controller struct {
  log *log.Entry

  udpConn *net.UDPConn    // listening on port
  udpAddr *net.UDPAddr    // sending to broadcast
}

func (controller *Controller) run() {
  var buf = make([]byte, MTU)

  if err := controller.sendPoll(); err != nil {
    controller.log.Fatalf("run: sendPoll: %v", err)
  }

  for {
    if read, udpAddr, err := controller.udpConn.ReadFromUDP(buf); err != nil {
      controller.log.Fatalf("run: UDPConn.ReadFromUDP: %v", err)
    } else if err := controller.recv(buf[:read]); err != nil {
      controller.log.Warnf("recv %v: %v", udpAddr, err)
    }
  }
}

func (controller *Controller) Run() {
  controller.run()
}

func (controller *Controller) Start() {
  go controller.run()
}

func (controller *Controller) recv(buf []byte) error {
  // decode header
  var header ArtHeader

  if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &header); err != nil {
    return fmt.Errorf("binary.Read ArtHeader: %v", err)
  }

  if header.ID != ARTNET {
    return fmt.Errorf("Invalid magic")
  }

  // decode packet
  var packet ArtPacket

  switch header.OpCode {
  case OpPoll:
    packet = &ArtPoll{}

  case OpPollReply:
    packet = &ArtPollReply{}

  case OpDmx:
    packet = &ArtDmx{}

  default:
    return fmt.Errorf("Unknown opcode: %04x", header.OpCode)
  }

  if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, packet); err != nil {
    return fmt.Errorf("binary.Read %T: %v", packet, err)
  }

  return controller.recvPacket(packet)
}

func (controller *Controller) recvPacket(packet ArtPacket) error {
  switch opPacket := packet.(type) {
  case *ArtPoll:
    if opPacket.ProtVer < ProtVer {
      return fmt.Errorf("Invalid protocol version: %v < %v", opPacket.ProtVer, ProtVer)
    }

    // ignore
    return nil

  case *ArtPollReply:
    return controller.recvPollReply(opPacket)

  case *ArtDmx:
    if opPacket.ProtVer < ProtVer {
      return fmt.Errorf("Invalid protocol version: %v < %v", opPacket.ProtVer, ProtVer)
    }

    // ignore
    return nil
  }

  return nil
}

func (controller *Controller) recvPollReply(pollReply *ArtPollReply) error {
  controller.log.Infof("recvPollReply: %#v", pollReply)

  return nil
}

// Send broadcast message
func (controller *Controller) send(packet ArtPacket) error {
  var buf bytes.Buffer

  if err := binary.Write(&buf, binary.LittleEndian, packet); err != nil {
    return err
  }

  controller.log.Debugf("send: opcode=%04x len=%v", packet.Header().OpCode, buf.Len())

  if _, err := controller.udpConn.WriteToUDP(buf.Bytes(), controller.udpAddr); err != nil {
    return err
  }

  return nil
}

func (controller *Controller) sendPoll() error {
  return controller.send(ArtPoll{
    ArtHeader: ArtHeader{
      ID: ARTNET,
      OpCode: OpPoll,
    },
    ProtVer: ProtVer,
  })
}
