import * as _ from 'lodash';
import { Observer } from 'rxjs/Observer';

import {
  DMX,
  Value,
  Color,
  Colors,
  ChannelType,
  HeadType,
  HeadConfig,
} from './types';

import {
  APIChannel,
  APIIntensity,
  APIColor,
  APIHead,
  APIGroup,
  APIPreset,
  APIParameters,
  APIChannelParameters,
  APIHeadParameters,
  PresetConfig
} from './api';

// POST API plumbing
export interface Post {
  type: String
  id: string
  parameters: Object
};

interface PostFunc {
  (parameters: APIParameters);
}
interface PostHeadFunc {
  (parameters: APIHeadParameters);
}

// Head.Intensity, Group.Intensity
export class IntensityParameter {
  private intensity: Value;

  constructor(private post: PostFunc, api: APIIntensity) {
    this.load(api)
  }
  load(api: APIIntensity) {
    this.intensity = api.Intensity;
  }

  get Intensity(): Value {
    return this.intensity;
  }
  set Intensity(value: Value) {
    this.post({Intensity: { Intensity: value } });
  }
}

// Head.Color, Group.Color
export class ColorParameter implements Color {
  red:        Value;
  green:      Value;
  blue:       Value;

  constructor(private post: PostFunc, api: APIColor) {
    this.load(api)
  }
  load(api: APIColor) {
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

  /* Post RGB values from Color */
  apply(color: Color) {
    this.post({Color: color});
  }
}

export class Channel {
  ID:       string;
  Type:     ChannelType;
  Address:  number;
  dmx:      DMX;
  value:    Value;

  constructor(private post: PostHeadFunc, api: APIChannel) {
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

// Common API for Head and Group
export interface Parameters {
  Intensity?: IntensityParameter;
  Color?:     ColorParameter;
  Colors?:    Colors;
}

export class Head implements Parameters {
  private post: PostHeadFunc;

  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;

  channels:   {[id: string]: Channel};
  Intensity?: IntensityParameter;
  Color?:     ColorParameter;

  constructor(postObserver: Observer<Post>, api: APIHead) {
    this.ID = api.ID;
    this.Type = api.Type;
    this.Config = api.Config;

    this.post = (headParameters: APIHeadParameters) => postObserver.next({type: "heads", id: this.ID, parameters: headParameters});
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
      this.Intensity = new IntensityParameter(this.post, api.Intensity);
    }
    if (api.Color) {
      this.Color = new ColorParameter(this.post, api.Color);
    }
  }

  /* Channel objects */
  get Channels(): Channel[] {
    let channels = Object.keys(this.channels).map(key => this.channels[key]);

    return _.sortBy(channels, channel => channel.Address);
  }

  get Colors(): Colors {
    return this.Type.Colors;
  }
}

export class Group implements Parameters {
  private post: PostFunc;

  ID:     string
  Heads:  Head[];

  Intensity?: IntensityParameter;
  Color?: ColorParameter;

  constructor(postObserver: Observer<Post>, api: APIGroup, heads: Head[]) {
    this.ID = api.ID;
    this.Heads = heads;

    this.post = (parameters: APIParameters) => postObserver.next({type: "groups", id: this.ID, parameters: parameters});
    this.load(api);
  }

  load(api: APIGroup) {
    console.log("Group.load", this.ID, api);

    if (api.Intensity) {
      this.Intensity = new IntensityParameter(this.post, api.Intensity);
    }
    if (api.Color) {
      this.Color = new ColorParameter(this.post, api.Color);
    }
  }

  // TODO: Colors
}

export class Preset {
  private post: PostFunc;

  ID:     string
  Config: PresetConfig

  constructor(postObserver: Observer<Post>, api: APIPreset) {
    this.ID = api.ID
    this.Config = api.Config

    this.post = (parameters: APIParameters) => postObserver.next({type: "presets", id: this.ID, parameters: parameters})
  }

  get Name(): string {
    if (this.Config.Name)
      return this.Config.Name
    else
      return this.ID
  }

  apply() {
    this.post({});
  }
}
