import { Component, Input, HostBinding } from '@angular/core';

import { Color } from './types';
import { Head } from './head';

@Component({
  moduleId: module.id,
  selector: 'dmx-head',

  templateUrl: 'head.component.html',
  styleUrls: [ 'head.component.css' ],
})
export class HeadComponent {
  @Input() head: Head;
}
