package artnet

type UniverseAddress struct {
  Net     uint8 // 0-128
  SubUni  uint8
}
type Universe []uint8


type ArtDmx struct {
  ArtHeader
  ProtVer uint16

  Sequence    uint8
  Physical    uint8
  SubUni      uint8
  Net         uint8
  Length      uint16

  Data        []uint8
}
