export interface ChannelType {
  Control?: string;
  Intensity?: boolean;
  Color?: string;
}
export interface HeadType {
  Vendor: string;
  Model:  string;
  Mode:   string;
  Channels: ChannelType[];
}
export interface HeadConfig {
  Type:     string;
  Universe: number;
  Address:  number;
}
export interface Channel {
  Type:     ChannelType;
  Address:  number;
  DMX:      number;
  Value:    number;
}
export class Head {
  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;
  Channels: Channel[];

  cmpAddress(other) : number {
    if (this.Config.Universe != other.Config.Universe)
      return this.Config.Universe - other.Config.Universe;

    if (this.Config.Address != other.Config.Address)
      return this.Config.Address - other.Config.Address;

    return 0;
  }
}
