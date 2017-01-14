import { Component } from '@angular/core';

import { Preset } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-presets',
  host: { class: 'view' },

  templateUrl: 'presets.component.html',
  styleUrls: [ 'presets.component.css' ],
})
export class PresetsComponent {
  constructor (private service: HeadService) {

  }

  list(): Preset[] {
    return this.service.listPresets();
  }

  isActive(preset: Preset): boolean {
    return false
  }

  click(preset: Preset) {
    preset.apply();
  }
}
