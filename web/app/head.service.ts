import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import { WebSocketService, WebSocketError } from 'websocket';

import * as _ from 'lodash';
import { Observable } from 'rxjs/Observable';
import { Observer } from 'rxjs/Observer';
import { Subscription } from 'rxjs/Subscription';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/map';

import { Head, Post, APIHead, APIEvents } from './head';

@Injectable()
export class HeadService {
  private webSocket: Subscription;
  postSubject = new Subject<Post>();
  heads: Map<string, Head>;
  active: Head = null;

  list(sort?: (Head) => any, filter?: (Head) => boolean) {
    let heads = Object.keys(this.heads).map(key => this.heads[key]);

    if (filter)
      heads = _.filter(heads, filter);

    if (sort)
      heads = _.sortBy(heads, sort);

    return heads;
  }
  byAddress(): Head[] {
    return this.list(head => [head.Config.Universe, head.Config.Address]);
  }
  byID(): Head[] {
    return this.list(head => head.ID);
  }
  byIntensity(): Head[] {
    return this.list(head => head.ID, head => !!head.Intensity);
  }

  select(head: Head) {
    console.log("Select head", head);
    this.active = head;
  }
  selected(head: Head): boolean {
    return this.active == head;
  }

  constructor(private http: Http, webSocketService: WebSocketService) {
    this.heads = {};

    this.get('/api/heads').subscribe(
      headsMap => {
        this.load(headsMap);

        console.log("Loaded heads", this.heads);
      }
    );

    this.postSubject.subscribe(
      headStream => {
        console.log(`POST head ${headStream.head.ID}...`, headStream.headPost);

        this.post(`/api/heads/${headStream.head.ID}`, headStream.headPost).subscribe(
          (headParams: APIHead) => {
            console.log(`POST head ${headStream.head.ID}: OK`, headParams);
          }
        );
      }
    );

    this.webSocket = webSocketService.connect<APIEvents>('/events').subscribe(
      (apiEvents: APIEvents) => {
        console.log("WebSocket APIEvents", apiEvents);

        this.load(apiEvents.Heads);
      },
      (error: WebSocketError) => {
        console.log("WebSocket Error", error);
      },
      () => {
        console.log("WebSocket Close");
      }
    );
  }
  private load(headsMap: Map<string, APIHead>) {
    for (let id in headsMap) {
      if (this.heads[id])
        this.heads[id].load(headsMap[id]);
      else
        this.heads[id] = new Head(this.postSubject, headsMap[id]);
    }
  }

  private get(url): any {
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
