import { Component, Input } from '@angular/core';

import { Head } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-intensity',
  host: { class: 'view dmx-heads' },
  templateUrl: 'intensity.component.html',
  styleUrls: [ 'intensity.component.css' ],
})
export class IntensityComponent {
  constructor (private service: HeadService) { }
}
