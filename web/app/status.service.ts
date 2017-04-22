import { Injectable } from '@angular/core';

import { Response } from '@angular/http';

@Injectable()
export class StatusService {
  // Managed by APIService from Websocket state
  connected: Boolean;

  error?: Error;

  requests_pending: number = 0;

  Connected() {
    this.connected = true;
  }
  Disconnected(error?: Error) {
    this.connected = false;

    if (error) {
      this.setError(error);
    }
  }

  TrackRequest(method: string, promise: Promise<any>) {
    this.requests_pending += 1;

    return promise.then(
      (value) => {
        this.requests_pending -= 1;

        return value;
      },
      (error) => {
        this.requests_pending -= 1;

        if (error instanceof Response) {
          this.setError(new Error(error.toString()));

        } else if (error instanceof Error) {
          this.setError(error);

        } else {
          this.setError(new Error(error));
        }

        throw error;
      }
    );
  }

  // Called by AppErrorHandler
  setError(error: Error) {
    this.error = error;
  }

  state(): string {
    if (!this.connected) {
      return 'offline';
    } else if (this.requests_pending > 0) {
      return 'upload';
    } else {
      return 'download';
    }
  }
}
