import { Component } from '@angular/core';

import { Preset, Group, Head } from './head';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-presets',
  host: { class: 'split view' },

  templateUrl: 'presets.component.html',
  styleUrls: [ 'presets.component.css' ],
})
export class PresetsComponent {
  preset?: Preset

  constructor (private api: APIService) {

  }

  list(): Preset[] {
    return this.api.listPresets();
  }

  isActive(preset: Preset): boolean {
    return preset == this.preset
  }

  click(preset: Preset) {
    preset.apply()
    this.preset = preset
  }

  listGroups(): Group[] {
    return this.preset.Groups ? Array.from(this.preset.Groups.keys()) : [];
  }
  listHeads(): Head[] {
    return this.preset.Heads ? Array.from(this.preset.Heads.keys()) : [];
  }
}
