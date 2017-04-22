import { Injectable } from '@angular/core';

@Injectable()
export class StatusService {
  // Set by the AppErrorHandler
  error?: Error;

  setError(error: Error) {
    this.error = error;
  }
}
