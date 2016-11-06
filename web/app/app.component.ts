import { Component, OnInit } from '@angular/core';

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
export class AppComponent implements OnInit {
  title = "qmsk::dmx";

  constructor (private headService: HeadService) { }

  ngOnInit(): void {
    this.headService.load().subscribe(
      heads => { console.log("Loaded heads", heads); }
    );
  }
}
