import { Component, Input } from '@angular/core';

import { Head, Channel } from './head';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-channels',
  host: { class: 'view' },
  templateUrl: 'channels.component.html',
  styleUrls: [ 'channels.component.css' ],
})
export class ChannelsComponent {
  constructor (private api: APIService) { }

  listHeads(): Head[] {
    return this.api.listHeads(head => [head.Config.Universe, head.Config.Address]);
  }

  setHeadChannelDMX(head: Head, channel: Channel, value: string) {
    channel.DMX = parseInt(value);
  }
}
