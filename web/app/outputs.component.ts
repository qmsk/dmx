import { Component, Input } from '@angular/core';

import { APIOutput } from './api';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-outputs',
  host: { class: 'view' },
  templateUrl: 'outputs.component.html',
})
export class OutputsComponent {
  constructor (private api: APIService) { }

  list(): APIOutput[] {
    return this.api.listOutputs((output: APIOutput) => output.Universe);
  }
}
