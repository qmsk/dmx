import { ErrorHandler, Injectable } from '@angular/core'

import { StatusService } from './status.service'

@Injectable()
export class AppErrorHandler implements ErrorHandler {
  constructor(private status: StatusService) { }

  handleError(error) {
    console.error(error);

    this.status.AppError(error);
  }
}
