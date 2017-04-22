import { Component, OnInit } from '@angular/core';

import { WebSocketService } from 'lib/websocket';

import { StatusService } from './status.service'
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-app',
  templateUrl: 'app.component.html',
  styleUrls: [ 'app.component.css' ],
  providers: [
    WebSocketService,
    APIService,
  ],
})
export class AppComponent implements OnInit {
  title = "qmsk::dmx";

  constructor (private api: APIService, private status: StatusService) { }

  ngOnInit() {
    console.log("App initialized");
    this.status.AppOK();
  }
}
