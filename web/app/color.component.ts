import { Component, Input, HostBinding } from '@angular/core';

import { Value, Color } from './types';
import { Head, APIColors } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-color',
  host: { class: 'view split' },

  templateUrl: 'color.component.html',
  styleUrls: [ 'color.component.css' ],
})
export class ColorComponent {
  colors: APIColors;
  color: Color;
  heads: Set<Head>;

  constructor (private service: HeadService) {
    this.heads = new Set<Head>();
  }

  headActive(head: Head): boolean {
    return this.heads.has(head);
  }

  colorActive(color: Color): boolean {
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

  /* Build new colors map from active heads */
  loadColors(): APIColors {
    // XXX: just return from first selected head..
    // TODO: merge color maps from multiple heads?
    for (let head of Array.from(this.heads)) {
      return head.Type.Colors;
    }
  }

  select(head: Head) {
    this.heads.add(head);

    this.colors = this.loadColors();
    this.color = this.loadColor(head.Color);
  }

  unselect(head: Head) {
    this.heads.delete(head);

    if (this.heads.size > 0) {
      this.colors = this.loadColors();
    } else {
      this.colors = this.color = null;
    }
  }

  apply(color: Color) {
    // XXX: this does not update the <dmx-head-colors>
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
