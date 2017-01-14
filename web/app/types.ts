export type DMX = number; // Integer 0 .. 255
export type Value = number; // Float 0.0 .. 1.0

export interface Color {
  Red:    Value;
  Green:  Value;
  Blue:   Value;
}
export type Colors = {[ID: string]: Color};

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
  Colors:   Colors;
}
export interface HeadConfig {
  Type:     string;
  Universe: number;
  Address:  number;
  Count:    number;
}
