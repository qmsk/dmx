import { Component, Input, HostBinding } from '@angular/core';

import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-color',
  host: { class: 'view flow' },

  templateUrl: 'color.component.html',
  styleUrls: [ 'color.component.css' ],
})
export class ColorComponent {
  constructor (private service: HeadService) { }
}
