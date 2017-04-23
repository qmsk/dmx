import { Component, Input } from '@angular/core';

import { Value } from './types';
import { Head } from './head';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-intensity',
  host: { class: 'view dmx-heads' },
  templateUrl: 'intensity.component.html',
  styleUrls: [ 'intensity.component.css' ],
})
export class IntensityComponent {
  constructor (private api: APIService) { }

  listHeads(): Head[] {
    return this.api.listHeads(head => head.ID, head => !!head.Intensity);
  }

  change(head: Head, value: Value) {
    head.Intensity.apply(value);
  }
}
