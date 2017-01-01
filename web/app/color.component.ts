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
  colors: Color[];
  color: Color;
  heads: Set<Head>;

  constructor (private service: HeadService) {
    this.heads = new Set<Head>();
  }

  hexByte(value: Value): string {
    return _.padStart(Math.trunc(value * 255).toString(16), 2, '0');
  }

  hexRGB(color: Color): string {
    return "#" + this.hexByte(color.Red) + this.hexByte(color.Green) + this.hexByte(color.Blue);
  }

  headActive(head: Head): boolean {
    return this.heads.has(head);
  }

  isActive(color: Color): boolean {
    return this.color
      && color.Red == this.color.Red
      && color.Green == this.color.Green
      && color.Blue == this.color.Blue
    ;
  }

  loadColor(color: Color): Color {
    return {
      Red:   color.Red,
      Green: color.Green,
      Blue:  color.Blue,
    };
  }

  /* Build new colors map from active heads
   * Optionally override any colors from given head.
   */
  loadColors(): Color[] {
    let colors = new Map<string, Color>();

    this.heads.forEach((head) => {
      for (let name in head.Type.Colors) {
        colors.set(name, head.Type.Colors[name]);
      }
    });

    return Array.from(colors.values());
  }

  select(head: Head) {
    this.heads.add(head);
    this.color = this.loadColor(head.Color);
    this.colors = this.loadColors();
  }

  unselect(head: Head) {
    this.heads.delete(head);
    // XXX: this.color = ...
    this.colors = this.loadColors();
  }

  apply(color: Color) {
    this.color = color;
    this.heads.forEach((head) => {
      head.Color.apply(color);
    });
  }

  /* Copy and apply color */
  click(color: Color) {
    this.apply(this.loadColor(color));
  }

  /* Change and apply color */
  change(channel: string, value: Value) {
    this.color[channel] = value;
    this.apply(this.color);
  }
}
