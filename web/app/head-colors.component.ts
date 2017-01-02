import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Value, Color } from './types';
import { APIColors } from './head';

@Component({
  moduleId: module.id,
  selector: 'dmx-head-colors',

  templateUrl: 'head-colors.component.html',
  styleUrls: [ 'head-colors.component.css' ],
})
export class HeadColorsComponent {
  @Input() colors: APIColors;
  @Input() color: Color;
  @Output() colorChange = new EventEmitter<Color>();

  makeColors(): Color[] {
    return Object.values(this.colors);
  }

  active(color: Color): boolean {
    return this.color
      && color.Red == this.color.Red
      && color.Green == this.color.Green
      && color.Blue == this.color.Blue
    ;
  }

  click(color: Color) {
    this.color = color;
    this.colorChange.emit(color);
  }
}
