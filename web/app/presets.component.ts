import { Component } from '@angular/core';

import { Value } from './types';
import { APIParameters } from './api';
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
  intensity: Value
  preset?: Preset

  constructor (private api: APIService) {
    this.intensity = 1.0
  }

  list(): Preset[] {
    return this.api.listPresets();
  }

  isActive(preset: Preset): boolean {
    return preset == this.preset
  }

  click(preset: Preset) {
    preset.apply(this.intensity)
    this.preset = preset
  }
  update(intensity) {
    this.intensity = intensity;

    if (this.preset) {
      this.preset.apply(intensity)
    }
  }

  listGroups(): Group[] {
    return this.preset.Groups ? Array.from(this.preset.Groups.keys()) : []
  }
  listHeads(): Head[] {
    return this.preset.Heads ? Array.from(this.preset.Heads.keys()) : []
  }

  groupParameters(group: Group): APIParameters {
    return this.preset.Groups.get(group)
  }
  headParameters(head: Head): APIParameters {
    return this.preset.Heads.get(head)
  }

}
