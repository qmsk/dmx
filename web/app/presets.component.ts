import { Component } from '@angular/core';

import { Preset } from './head';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-presets',
  host: { class: 'view' },

  templateUrl: 'presets.component.html',
  styleUrls: [ 'presets.component.css' ],
})
export class PresetsComponent {
  constructor (private api: APIService) {

  }

  list(): Preset[] {
    return this.api.listPresets();
  }

  isActive(preset: Preset): boolean {
    return false
  }

  click(preset: Preset) {
    preset.apply();
  }
}
