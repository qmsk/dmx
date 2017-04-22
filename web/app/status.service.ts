import { Injectable } from '@angular/core';

@Injectable()
export class StatusService {
  // Managed by APIService from Websocket state
  connected: Boolean;

  error?: Error;

  Connected() {
    this.connected = true;
  }
  Disconnected(error?: Error) {
    this.connected = false;

    if (error) {
      this.setError(error);
    }
  }

  // Called by AppErrorHandler
  setError(error: Error) {
    this.error = error;
  }
}
