import { Component, OnInit } from '@angular/core';

import { Head } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-dashboard',
  host: { class: 'view split' },
  templateUrl: 'dashboard.component.html',
})
export class DashboardComponent {
  constructor (private headService: HeadService) { }
}
