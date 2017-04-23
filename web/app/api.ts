import { DMX, Value, Color, Colors, ChannelType, HeadType, HeadConfig } from './types';

// GET
export interface APIOutputArtnet {
  Name: string
  Description: string
  Ethernet: string
  Version: number
}

export interface APIOutput {
  Universe:   number
  Connected:  string

  Address:    string
  Port:       number
  Seen:       string
  Artnet:     APIOutputArtnet
}
export type APIOutputs = {[id: string]: APIOutput};

export interface APIChannel {
  ID:       string;
  Type:     ChannelType;
  Address:  number;
  DMX:      DMX;
  Value:    Value;
}

export interface APIIntensity {
  Intensity: Value;
}

export interface APIColor extends Color {

}

export interface PresetConfig {
  Name:     string;
  All?:   APIParameters
  Groups: {[id: string]: APIParameters}
  Heads:  {[id: string]: APIParameters}
}

export interface APIHead {
  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;

  Channels:   {[id: string]: APIChannel};
  Intensity?: APIIntensity;
  Color?:     APIColor;
}
export type APIHeads = {[id: string]: APIHead};

export interface APIGroup {
  ID:     string
  Heads:  string[]
  Colors: Colors

  Intensity?: APIIntensity
  Color?:     APIColor
}
export type APIGroups = {[id: string]: APIGroup}


export interface APIPreset {
  ID:   string
  Config: PresetConfig
}
export type APIPresets = {[id: string]: APIPreset}

export interface API {
  Outputs: APIOutputs
  Heads:   APIHeads
  Groups:  APIGroups
  Presets: APIPresets
}

// POST
export interface APIChannelParameters {
  DMX?: DMX;
  Value?: Value;
}

export interface APIParameters {
  Intensity?: APIIntensity;
  Color?:     APIColor;
}

export interface APIHeadParameters extends APIParameters {
  Channels?:  {[ID: string]: APIChannelParameters};
};

export interface APIPresetParameters {
  Intensity?: Value
}

// WebSocket
export interface APIEvents extends API {

}
