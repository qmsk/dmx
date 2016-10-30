import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import { Head } from './head';

@Injectable()
export class HeadService {
  constructor(private http: Http) {

  }

  list(): Observable<Head[]> {
    return this.http.get('/api/heads/')
      .map(response => response.json() as Head[])
    ;
  }
}
