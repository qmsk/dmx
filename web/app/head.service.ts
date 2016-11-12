import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import * as _ from 'lodash';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/map';

import { Head, ValueStream, HeadStream, PostFunc, Channel, HeadIntensity, HeadColor } from './head';

@Injectable()
export class HeadService {
  postSubject = new Subject<HeadStream>();
  heads: Head[];
  active: Head = null;

  load(headsData: Object[]) {
    this.heads = headsData.map(headData => new Head(this.postSubject, headData));
  }

  byAddress(): Head[] {
    return _.sortBy(this.heads, head => [head.Config.Universe, head.Config.Address]);
  }
  byID(): Head[] {
    return _.sortBy(this.heads, head => head.ID);
  }

  select(head: Head) {
    console.log("Select head", head);
    this.active = head;
  }
  selected(head: Head): boolean {
    return this.active == head;
  }

  constructor(private http: Http) {
    this.get('/api/heads/').subscribe(
      headsData => {
        this.load(headsData);

        console.log("Loaded heads", this.heads);
      }
    );

    this.postSubject.subscribe(
      headStream => {
        console.log(`Post head=${headStream.head.ID}`, headStream.valueStream);

        this.post(`/api/heads/${headStream.head.ID}`, headStream.valueStream).subscribe(
          headParams => {
            console.log(`Load head=${headStream.head.ID}`, headParams);

            headStream.head.load(headParams);
          }
        );
      }
    );
  }

  get(url): any {
    return this.http.get(url)
      .map(response => response.json())
    ;
  }

  private post(url, params): Observable<Object> {
    let headers = new Headers({
      'Content-Type': 'application/json',
    });
    let options = new RequestOptions({ headers: headers });

    return this.http.post(url, params, options)
      .map(response => response.json())
    ;
  }
}
