package artnet

type ArtDmx struct {
  ArtHeader
  ProtVer uint16

  Sequence    uint8
  Physical    uint8
  SubUni      uint8
  Net         uint8
  Length      uint16

  Data        []byte
}
