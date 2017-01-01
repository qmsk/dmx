import { Component, Input, HostBinding } from '@angular/core';

import * as _ from 'lodash';
import { Value, Color } from './types';
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
  activeColor: Color;
  colors: Color[];

  constructor (private service: HeadService) { }

  hexByte(value: Value): string {
    return _.padStart(Math.trunc(value * 255).toString(16), 2, '0');
  }

  hexRGB(color: Color): string {
    return "#" + this.hexByte(color.Red) + this.hexByte(color.Green) + this.hexByte(color.Blue);
  }

  isActive(color: Color): boolean {
    return this.activeColor
      && color.Red == this.activeColor.Red
      && color.Green == this.activeColor.Green
      && color.Blue == this.activeColor.Blue
    ;
  }

  select(head: Head) {
    this.service.select(head);
    this.activeColor = head.Color;
    this.colors = Object.values(head.Type.Colors);
  }

  apply(color: Color) {
    this.activeColor = color;
    this.service.apply(head => head.Color.apply(color));
  }
}
