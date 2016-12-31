import { Component, Input } from '@angular/core';

import { Head, Channel } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-channels',
  host: { class: 'view' },
  templateUrl: 'channels.component.html',
  styleUrls: [ 'channels.component.css' ],
})
export class ChannelsComponent {
  constructor (private headService: HeadService) { }

  heads() {
    return this.headService.byAddress();
  }

  setHeadChannelDMX(head: Head, channel: Channel, value: string) {
    channel.DMX = parseInt(value);
  }
}