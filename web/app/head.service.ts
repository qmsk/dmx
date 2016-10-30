import { Injectable } from '@angular/core';
import { Headers, Http } from '@angular/http';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import { Head } from './head';

@Injectable()
export class HeadService {
  constructor(private http: Http) {

  }

  private decode(headData) :Head {
    return Object.assign(new Head, headData);
  }

  list(): Observable<Head[]> {
    return this.http.get('/api/heads/')
      .map(response => response.json().map(headData => this.decode(headData)))
    ;
  }
}
