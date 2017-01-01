import { Component, Input, HostBinding } from '@angular/core';

import { Head } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-color',
  host: { class: 'view split' },

  templateUrl: 'color.component.html',
  styleUrls: [ 'color.component.css' ],
})
export class ColorComponent {
  color: string;

  constructor (private service: HeadService) { }

  select(head: Head) {
    this.service.select(head);
    this.color = head.Color.hexRGB();
  }

  setColor(color: string) {
    console.log("setColor", color);
  }
}
