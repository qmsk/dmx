import { Component, Input, HostBinding } from '@angular/core';

import { Color } from './types';
import { Parameters } from './head';

@Component({
  moduleId: module.id,
  selector: 'dmx-controls',
  host: { class: 'controls' },

  templateUrl: 'controls.component.html',
  styleUrls: [ 'controls.component.css' ],
})
export class ControlsComponent {
  @Input() parameters: Parameters;
}
