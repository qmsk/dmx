import { Component, Input } from '@angular/core';

import { Head } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-heads',
  host: { class: 'view split' },

  templateUrl: 'heads.component.html',
  styleUrls: [ 'heads.component.css' ],
})
export class HeadsComponent {
  constructor (private service: HeadService) { }
}
