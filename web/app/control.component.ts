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
    this.value = parseFloat(eventValue);

    console.log(`Update control ${this.label}`, this.value);
    
    this.valueChange.emit(this.value);
  }
}
