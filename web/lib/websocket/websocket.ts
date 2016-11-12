import { Injectable } from '@angular/core';

import { Subject, Observable, Observer } from 'rxjs/Rx';

export enum ReadyState {
    Connecting  = 0,
    Open        = 1,
    Closing     = 2,
    Closed      = 3,
}

export interface Error {
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

    };
    ws.onerror = (event: Event) => {
      console.log("WebSocket onerror", event);
    };
    ws.onmessage = (event: MessageEvent) => {
      subject.next(JSON.parse(event.data));
    };
    ws.onclose = (closeEvent: CloseEvent) => {
      let error = <Error>{code: closeEvent.code, reason: closeEvent.reason};

      if (closeEvent.wasClean) {
        subject.complete();
      } else {
        subject.error(error);
      }
    };

    return subject;
  }
}
