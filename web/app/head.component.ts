import { Component, Input, HostBinding } from '@angular/core';

import { Head } from './head';

@Component({
  moduleId: module.id,
  selector: 'dmx-head',

  templateUrl: 'head.component.html',
  styleUrls: [ 'head.component.css' ],
})
export class HeadComponent {
  @Input() head: Head;

  @HostBinding('class.active') get active() {
    return this.head ? this.head.active : false;
  }
}
