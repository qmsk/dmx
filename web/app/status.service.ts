import { Injectable } from '@angular/core';

import { Response } from '@angular/http';

export class Status {
  constructor(public icon: string, public message: string, public error?: Error) {

  }
}

@Injectable()
export class StatusService {
  app?: Status;
  request?: Status;
  websocket?: Status;

  // APIService WebSocket connection
  websocket_connected: Boolean;

  // APIService GET/POST requests
  requests_pending: number = 0;

  constructor() {
    this.AppLoading();
  }

  all(): Status[] {
    return [this.app, this.request, this.websocket].filter(status => !!status);
  }

  AppLoading() {
    this.app = new Status('autorenew', "Application loading...")
  }
  AppOK() {
    this.app = null;
  }
  AppError(error: Error) {
    this.app = new Status('error', error.toString(), error);
  }

  Connecting() {
    // keep disconnect error
    let error = this.websocket ? this.websocket.error : null;

    this.websocket_connected = false;
    this.websocket = new Status('cloud_queue', "Websocket Connecting...", error);
  }
  Connected() {
    // debounce, this can get called multiple times
    if (!this.websocket_connected) {
      this.websocket_connected = true;
      this.websocket = new Status('cloud', "Websocket Connected");
    }
  }
  Disconnected(error?: Error) {
    this.websocket_connected = false;

    if (error) {
      this.websocket = new Status('cloud_off', "Websocket Error", error);
    } else {
      this.websocket = new Status('cloud_off', "Websocket Disconnected");
    }
  }

  RequestOK() {
    this.request = null;
  }
  RequestError(error: Error) {
    this.request = new Status('warning', error.toString(), error);
  }
  TrackRequest(method: string, promise: Promise<any>) {
    this.requests_pending += 1;

    return promise.then(
      (value) => {
        this.requests_pending -= 1;

        this.RequestOK();

        return value;
      },
      (error) => {
        this.requests_pending -= 1;

        // XXX: this is horrible in here
        if (error instanceof Response) {
          this.RequestError(new Error(error.toString()));

        } else if (error instanceof Error) {
          this.RequestError(error);

        } else {
          this.RequestError(new Error(error));
        }

        throw error;
      }
    );
  }
}
