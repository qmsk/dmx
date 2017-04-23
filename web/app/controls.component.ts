import { Component, Input, HostBinding } from '@angular/core';

import { Color } from './types';
import { Parameters } from './head';

@Component({
  moduleId: module.id,
  selector: 'dmx-controls',
  host: { class: 'controls' },

  templateUrl: 'controls.component.html',
})
export class ControlsComponent {
  @Input() parameters: Parameters;
}
