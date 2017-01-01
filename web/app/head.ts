import * as _ from 'lodash';

import { DMX, Value } from './types';
import { Observer } from 'rxjs/Observer';

export interface ChannelPost {
  DMX?: DMX;
  Value?: Value;
}
export interface HeadPost {
  Channels?:   Map<string, ChannelPost>;
  Intensity?: APIHeadIntensity;
  Color?:     APIHeadColor;
};
export type Post = { head: Head, headPost: HeadPost };

interface HeadPostFunc {
  (post: HeadPost);
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

export interface APIHeadIntensity {
  Intensity: Value;
}
export class HeadIntensity {
  private intensity: Value;

  constructor(private post: HeadPostFunc, api: APIHeadIntensity) {
    this.load(api)
  }
  load(api: APIHeadIntensity) {
    this.intensity = api.Intensity;
  }

  get Intensity(): Value {
    return this.intensity;
  }
  set Intensity(value: Value) {
    this.post({Intensity: { Intensity: value } });
  }
}

export interface APIHeadColor {
  Red:    Value;
  Green:  Value;
  Blue:   Value;
}
export class HeadColor {
  red:        Value;
  green:      Value;
  blue:       Value;

  constructor(private post: HeadPostFunc, api: APIHeadColor) {
    this.load(api)
  }
  load(api: APIHeadColor) {
    this.red = api.Red;
    this.green = api.Green;
    this.blue = api.Blue;
  }

  get Red(): Value { return this.red; }
  get Green(): Value { return this.green; }
  get Blue(): Value { return this.blue; }

  set Red(value: Value) {
    this.post({Color: {
      Red: value,
      Green: this.green,
      Blue: this.blue,
    }});
  }
  set Green(value: Value) {
    this.post({Color: {
      Red: this.red,
      Green: value,
      Blue: this.blue,
    }});
  }
  set Blue(value: Value) {
    this.post({Color: {
      Red: this.red,
      Green: this.green,
      Blue: value,
    }});
  }

  hexField(value: Value): string {
    return _.padStart(Math.trunc(value * 255).toString(16), 2, '0');
  }

  hexRGB(): string {
    let color = "#" + this.hexField(this.Red) + this.hexField(this.Green) + this.hexField(this.Blue);

    console.log("Head.hexRGB", color);

    return color;
  }
}

export interface APIChannel {
  ID:       string;
  Type:     ChannelType;
  Address:  number;
  DMX:      DMX;
  Value:    Value;
}
export class Channel {
  ID:       string;
  Type:     ChannelType;
  Address:  number;
  dmx:      DMX;
  value:    Value;

  constructor(private post: HeadPostFunc, api: APIChannel) {
    this.ID = api.ID;
    this.Type = api.Type;
    this.Address = api.Address;

    this.load(api);
  }
  load(api: APIChannel) {
    console.log("\tChannel.load", this.ID, api);

    this.dmx = api.DMX;
    this.value = api.Value;
  }

  get DMX(): DMX { return this.dmx; }
  set DMX(value: DMX) {
    this.post({Channels: {[this.ID]: {DMX: value}}});
  }

  get Value(): Value { return this.value; }
  set Value(value: Value) {
    this.post({Channels: {[this.ID]: {Value: value}}});
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

export interface APIHead {
  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;

  Channels:   {[id: string]: APIChannel};
  Intensity?: APIHeadIntensity;
  Color?:     APIHeadColor;
}
export class Head {
  private post: HeadPostFunc;

  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;

  channels:   {[id: string]: Channel};
  Intensity?: HeadIntensity;
  Color?:     HeadColor;

  constructor(postObserver: Observer<Post>, api: APIHead) {
    this.ID = api.ID;
    this.Type = api.Type;
    this.Config = api.Config;

    this.post = (headPost: HeadPost) => postObserver.next({head: this, headPost: headPost});
    this.channels = {};
    this.load(api);
  }
  load(api: APIHead) {
    console.log("Head.load", this.ID, api);

    if (api.Channels) {
      for (let channelID in api.Channels) {
        let channel = this.channels[channelID]; if (channel) {
          channel.load(api.Channels[channelID]);
        } else {
          this.channels[channelID] = new Channel(this.post, api.Channels[channelID]);
        }
      }
    }
    if (api.Intensity) {
      this.Intensity = new HeadIntensity(this.post, api.Intensity);
    }
    if (api.Color) {
      this.Color = new HeadColor(this.post, api.Color);
    }
  }

  /* Channel objects */
  get Channels(): Channel[] {
    let channels = Object.keys(this.channels).map(key => this.channels[key]);

    return _.sortBy(channels, channel => channel.Address);
  }
}

export interface APIEvents {
  Heads: Map<string, APIHead>;
}
