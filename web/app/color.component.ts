import { Component, Input, HostBinding } from '@angular/core';

import { Value, Color, Colors } from './types';
import { Head, Group, ColorParameter } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-color',
  host: { class: 'view split' },

  templateUrl: 'color.component.html',
  styleUrls: [ 'color.component.css' ],
})
export class ColorComponent {
  colors: Colors;
  color: Color;
  heads: Set<Head>;
  groups: Set<Group>;

  constructor (private service: HeadService) {
    this.heads = new Set<Head>();
    this.groups = new Set<Group>();
  }

  listHeads(): Head[] {
    return this.service.listHeads(head => head.ID, head => !!head.Color);
  }
  listGroups(): Group[] {
    return this.service.listGroups(group => group.ID, group => !!group.Color);
  }

  headActive(head: Head): boolean {
    return this.heads.has(head);
  }
  groupActive(group: Group): boolean {
    return this.groups.has(group);
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
  loadColors(): Colors {
    // XXX: just return from first selected head or group...
    // TODO: merge color maps from multiple heads?
    for (let head of Array.from(this.heads)) {
      return head.Type.Colors;
    }

    for (let group of Array.from(this.groups)) {
      for (let head of group.Heads) {
        return head.Type.Colors;
      }
    }

    return null;
  }

  selectHead(head: Head) {
    this.heads.add(head);

    this.colors = this.loadColors();
    this.color = this.loadColor(head.Color);
  }
  selectGroup(group: Group) {
    this.groups.add(group);

    this.colors = this.loadColors();
    this.color = this.loadColor(group.Color);
  }

  unselectHead(head: Head) {
    this.heads.delete(head);

    if (this.heads.size > 0) {
      this.colors = this.loadColors();
    } else {
      this.colors = this.color = null;
    }
  }
  unselectGroup(group: Group) {
    this.groups.delete(group);

    if (this.groups.size > 0) {
      this.colors = this.loadColors();
    } else {
      this.colors = this.color = null;
    }
  }

  apply(color: Color) {
    // XXX: this does not update the <dmx-head-colors> component's color Input?
    this.color = color;
    this.heads.forEach((head) => {
      head.Color.apply(color);
    });
    this.groups.forEach((group) => {
      group.Color.apply(color);
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
