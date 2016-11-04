import { Component, Input, HostBinding } from '@angular/core';

import { Value, DMX } from './types';

@Component({
  moduleId: module.id,
  selector: 'dmx-control',

  templateUrl: 'control.component.html',
  styleUrls: [ 'control.component.css' ],
})
export class ControlComponent {
  @Input() label: string;
  @Input() value: Value;
}
