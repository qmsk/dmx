import * as _ from 'lodash';

import { DMX, Value } from './types';
import { Observer } from 'rxjs/Observer';

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
  private post: PostFunc;

  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;

  channels = new Map<string, Channel>();
  Intensity?: HeadIntensity;
  Color?:     HeadColor;

  constructor(postObserver: Observer<HeadStream>, data: Object) {
    this.ID = data['ID'];
    this.Type = data['Type'];
    this.Config = data['Config'];

    this.post = (valueStream: ValueStream) => postObserver.next({head: this, valueStream: valueStream});
    this.load(data);
  }
  load(data: Object) {
    let channelsData = data['Channels']; if (channelsData) {
      for (let channelID in channelsData) {
        let channel = this.channels[channelID]; if (channel) {
          this.channels[channelID].load(channelsData[channelID]);
        } else {
          this.channels[channelID] = new Channel(this.post, channelsData[channelID]);
        }
      }
    }
    let intensityData = data['Intensity']; if (intensityData) {
      this.Intensity = new HeadIntensity(this.post, intensityData);
    }
    let colorData = data['Color']; if (colorData) {
      this.Color = new HeadColor(this.post, colorData);
    }
  }

  /* Channel objects */
  get Channels(): Channel[] {
    let channels = Object.keys(this.channels).map(key => this.channels[key]);

    return _.sortBy(channels, channel => channel.Address);
  }
}
