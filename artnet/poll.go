package artnet

type ArtPoll struct {
  ArtHeader
  ProtVer uint16

  TalkToMe    uint8
  Priority    uint8
}

type ArtPollReply struct {
  ArtHeader

  IPAddress   [4]byte
  PortNumber  uint16
  VersInfo    uint16
  NetSwitch   uint8
  SubSwitch   uint8
  Oem         uint16
  UbeaVersion uint8
  Status1     uint8
  EstaMan     uint16
  ShortName   [18]byte
  LongName    [64]byte
  NodeReport  [64]byte
  NumPorts    uint16
  PortTypes   [4]uint8
  GoodInput   [4]uint8
  GoodOutput  [4]uint8
  SwIn        [4]uint8
  SwOut       [4]uint8
  SwVideo     uint8
  SwMacro     uint8
  SwRemote    uint8
  Spare1      byte
  Spare2      byte
  Spare3      byte
  Style       byte
  Mac         [6]byte
  BindIp      [4]byte
  BindIndex   uint8
  Status2     uint8
}
