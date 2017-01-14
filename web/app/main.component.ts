import { Component, Input } from '@angular/core';

import { Head, Group } from './head';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-main',
  host: { class: 'view split' },

  templateUrl: 'main.component.html',
  styleUrls: [ 'main.component.css' ],
})
export class MainComponent {
  head: Head
  group: Group

  constructor (private api: APIService) {

  }

  listHeads(): Head[] {
    return this.api.listHeads(head => head.ID)
  }
  listGroups(): Group[] {
    return this.api.listGroups(group => group.ID)
  }

  selectHead(head: Head) {
    this.head = head
    this.group = null
  }
  selectGroup(group: Group) {
    this.head = null
    this.group = group
  }

  headActive(head: Head): boolean {
    return this.head == head
  }
  headSemiActive(head: Head): boolean {
    if (this.group) {
      return this.group.Heads.has(head)
    } else {
      return false
    }
  }
  groupActive(group: Group): boolean {
    return this.group == group;
  }
}
