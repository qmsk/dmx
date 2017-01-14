import { Component } from '@angular/core';

import { WebSocketService } from 'lib/websocket';

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
export class AppComponent  {
  title = "qmsk::dmx";

  constructor (private api: APIService) { }
}
