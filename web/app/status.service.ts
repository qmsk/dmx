import { Injectable } from '@angular/core';

export class Status {
  constructor(public icon: string, public message: string, public error?: Error) {

  }
}

@Injectable()
export class StatusService {
  app?: Status;
  requests?: Status;
  websocket?: Status;

  // APIService WebSocket connection
  websocket_connected: Boolean;

  // APIService GET/POST requests
  requests_pending: number = 0;

  constructor() {
    this.AppLoading();
  }

  all(): Status[] {
    return [this.app, this.requests, this.websocket].filter(status => !!status);
  }
  status(): Status {
    if (this.requests_pending > 0) {
      return this.requests;
    } else {
      return this.websocket;
    }
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

  RequestStart(method: string, url: string) {
    this.requests_pending += 1;

    this.requests = new Status(method == 'GET' ? 'cloud_download' : 'cloud_upload', method + " " + url);
  }
  RequestEnd(method: string, url: string, error?: Error) {
    this.requests_pending -= 1;

    if (error) {
      this.requests = new Status('warning', method + " " + url, error);
    } else if (this.requests_pending == 0){
      this.requests = null;
    }
  }
}
