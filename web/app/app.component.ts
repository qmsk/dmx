import { Component } from '@angular/core';

import { HeadService } from './head.service';
import { Head } from './head';

@Component({
  moduleId: module.id,
  selector: 'dmx-app',
  templateUrl: 'app.component.html',
  styleUrls: [ 'app.component.css' ],
  providers: [
    HeadService,
  ],
})
export class AppComponent  {
  title = "qmsk::dmx";

  constructor (private headService: HeadService) { }
}
