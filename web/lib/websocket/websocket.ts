import { Injectable } from '@angular/core';

import { Observable } from 'rxjs/Observable';
import { Observer } from 'rxjs/Observer';

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
  public connect<T>(url: string): Observable<T> {
    if (url[0] == '/') {
      url = (location.protocol == 'http:' ? 'ws:' : 'wss:') + '//' + location.host + url;
    }

    return Observable.create(function (observer: Observer<T>) {
      console.log("WebSocket connect");

      let ws = new WebSocket(url);

      ws.onopen = (event: Event) => {
        // wait for first message before advancing subject
        console.log("WebSocket open");
      };
      ws.onerror = (event: Event) => {
        // the websocket error event is intentionally useless, it does not contain any other information
        // just wait for the close event...
        console.log("WebSocket error", event);
      };
      ws.onmessage = (event: MessageEvent) => {
        observer.next(JSON.parse(event.data));
      };
      ws.onclose = (closeEvent: CloseEvent) => {
        console.log("WebSocket close", closeEvent);

        let error = new WebSocketError("Websocket closed with code=" + closeEvent.code + ": " + closeEvent.reason);
        error.code = closeEvent.code;
        error.reason = closeEvent.reason;

        if (closeEvent.wasClean) {
          observer.complete();
        } else {
          observer.error(error);
        }
      };

      return function() {
        console.log("WebSocket disconnect");

        ws.close();
      }
    });
  }
}
