import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Value, Color, Colors } from './types';

@Component({
  moduleId: module.id,
  selector: 'dmx-color-controls',

  templateUrl: 'color-controls.component.html',
  //styleUrls: [ 'color-controls.component.css' ],
})
export class ColorControlsComponent {
  @Input() color: Color;
  @Input() colors: Colors;
  @Output() colorChange = new EventEmitter<Color>();

  makeColors(): Color[] {
    return Object.keys(this.colors).map((id) => this.colors[id]);
  }

  active(color: Color): boolean {
    return (
         color.Red == this.color.Red
      && color.Green == this.color.Green
      && color.Blue == this.color.Blue
    );
  }

  load(color: Color) {
    return {
      Red:   color.Red,
      Green: color.Green,
      Blue:  color.Blue,
    };
  }

  /*
   * XXX: This replaces the external parameter state with our local control state... any external changes to the
   *      color state will not update the controls until this component is reloaded.
   *
   * TODO: Separate the remote parameter state and our local control state... the local state is needed to
   *       merge multiple pending changes, but once the changes have been applied, we should update back
   *       to the external parameter state...
   */
  apply(color: Color) {
    this.color = color;
    this.colorChange.emit(color);
  }

  /* Select and apply color*/
  click(color: Color) {
    this.apply(this.load(color));
  }

  /* Change and apply color */
  change(channel: string, value: Value) {
    let color = this.load(this.color);
    color[channel] = value;
    this.apply(color);
  }
}
