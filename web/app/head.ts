import { DMX, Value } from './types';

export type ValueStream = { [key: string]: ValueStream | Value };
export type HeadStream = { head: Head, valueStream: ValueStream };

export interface PostFunc {
  (stream: ValueStream);
}

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

export class HeadIntensity {
  private intensity: Value;

  constructor(private post: PostFunc, data: Object) {
    this.intensity = data['Intensity'];
  }

  get Intensity(): Value {
    return this.intensity;
  }
  set Intensity(value: Value) {
    this.post({"Intensity": { "Intensity": value } });
  }
}

export class HeadColor {
  red:        Value;
  green:      Value;
  blue:       Value;

  constructor(private post: PostFunc, data: Object) {
    this.load(data);
  }

  load(data: Object) {
    this.red = data['Red'];
    this.green = data['Green'];
    this.blue = data['Blue'];
  }

  get Red(): Value { return this.red; }
  get Green(): Value { return this.green; }
  get Blue(): Value { return this.blue; }

  set Red(value: Value) {
    this.post({"Color": {
      "Red": value,
      "Green": this.green,
      "Blue": this.blue,
    }});
  }
  set Green(value: Value) {
    this.post({"Color": {
      "Red": this.red,
      "Green": value,
      "Blue": this.blue,
    }});
  }
  set Blue(value: Value) {
    this.post({"Color": {
      "Red": this.red,
      "Green": this.green,
      "Blue": value,
    }});
  }
}

export class Channel {
  ID:       string;
  Type:     ChannelType;
  Address:  number;
  dmx:      DMX;
  value:    Value;

  constructor(private post: PostFunc, data: Object) {
    this.ID = data['ID'];
    this.Type = data['Type'];
    this.Address = data['Address'];

    this.load(data);
  }
  load(data: Object) {
    this.dmx = data['DMX'];
    this.value = data['Value'];
  }

  get DMX(): DMX { return this.dmx; }
  set DMX(value: DMX) {
    let channels = {};
    channels[this.ID] = { "DMX": value };
    this.post({"Channels": channels});
  }

  get Value(): Value { return this.value; }
  set Value(value: Value) {
    let channels = {};
    channels[this.ID] = { "Value": value },
    this.post({"Channels": channels});
  }

  typeClass(): string {
    if (this.Type.Control) {
      return "Control";
    } else if (this.Type.Intensity) {
      return "Intensity";
    } else if (this.Type.Color) {
      return "Color";
    } else {
      return "Unknown";
    }
  }
  typeLabel(): string {
    if (this.Type.Control) {
      return this.Type.Control;
    } else if (this.Type.Intensity) {
      return "Intensity";
    } else if (this.Type.Color) {
      return this.Type.Color;
    } else {
      return "Unknown";
    }
  }
}
export class Head {
  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;

  channels = new Map<string, Channel>();
  Intensity?: HeadIntensity;
  Color?:     HeadColor;

  constructor(data: Object) {
    this.ID = data['ID'];
    this.Type = data['Type'];
    this.Config = data['Config'];
  }

  /* Channel objects */
  get Channels(): Channel[] {
    return Object.keys(this.channels).map(key => this.channels[key]);
  }

  cmpHead(other): number {
    if (this.ID < other.ID)
      return -1;
    else if (this.ID > other.ID)
      return +1;
    else
      return 0;
  }

  cmpAddress(other): number {
    if (this.Config.Universe != other.Config.Universe)
      return this.Config.Universe - other.Config.Universe;

    if (this.Config.Address != other.Config.Address)
      return this.Config.Address - other.Config.Address;

    return 0;
  }
}
