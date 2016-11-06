import { Component, Input, Output, EventEmitter } from '@angular/core';

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
  @Output() valueChange = new EventEmitter<Value>();

  update(eventValue :string) {
    let value = parseFloat(eventValue);

    console.log(`Update control ${this.label}`, value);

    this.valueChange.emit(value);
  }
}
