import { Component, Input } from '@angular/core';

import { APIParameters } from './api';
import { Preset } from './head';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-preset-parameters',

  templateUrl: 'preset-parameters.component.html',
  styleUrls: [ 'preset-parameters.component.css' ],

})
export class PresetParametersComponent {
  @Input() title: string
  @Input() parameters: APIParameters
}
