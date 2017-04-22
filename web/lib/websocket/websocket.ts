import { Injectable } from '@angular/core';

import { Observable } from 'rxjs/Observable';
import { Observer } from 'rxjs/Observer';
import { Subject } from 'rxjs/Subject';

export enum ReadyState {
    Connecting  = 0,
    Open        = 1,
    Closing     = 2,
    Closed      = 3,
}

export class WebSocketError extends Error {
  code: number;
  reason: string;
}

@Injectable()
export class WebSocketService {
  public connect<T>(url: string): Subject<T> {
    if (url[0] == '/') {
      url = (location.protocol == 'http:' ? 'ws:' : 'wss:') + '//' + location.host + url;
    }
    let ws = new WebSocket(url);
    let subject = new Subject<T>();

    ws.onopen = (event: Event) => {
      // wait for first message before advancing subject
    };
    ws.onerror = (event: Event) => {
      // the websocket error event is intentionally useless, it does not contain any other information
      // just wait for the close event...
      console.log("WebSocket onerror", event);
    };
    ws.onmessage = (event: MessageEvent) => {
      subject.next(JSON.parse(event.data));
    };
    ws.onclose = (closeEvent: CloseEvent) => {
      let error = new WebSocketError("Websocket closed with code=" + closeEvent.code + ": " + closeEvent.reason);
      error.code = closeEvent.code;
      error.reason = closeEvent.reason;

      if (closeEvent.wasClean) {
        subject.complete();
      } else {
        subject.error(error);
      }
    };

    return subject;
  }
}
