import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import { WebSocketService, WebSocketError } from 'lib/websocket';

import * as _ from 'lodash';
import { Observable } from 'rxjs/Observable';
import { Observer } from 'rxjs/Observer';
import { Subscription } from 'rxjs/Subscription';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/map';

import { API, APIEvents, APIHeads, APIHeadParameters } from './api';
import { Head, Post } from './head';

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
  ofIntensity(): Head[] {
    return this.list(head => head.ID, head => !!head.Intensity);
  }
  ofColor(): Head[] {
    return this.list(head => head.ID, head => !!head.Color);
  }

  select(head: Head) {
    console.log("Select head", head);
    this.active = head;
  }
  isActive(head: Head): boolean {
    return this.active == head;
  }
  apply(func: (Head) => any) {
    if (this.active)
      func(this.active);
  }

  constructor(private http: Http, webSocketService: WebSocketService) {
    this.heads = new Map<string, Head>();

    this.get('/api/').subscribe(
      api => {
        this.loadHeads(api.Heads);
      }
    );

    this.postSubject.subscribe(
      post => {
        console.log(`POST head ${post.head.ID}...`, post.headParameters);

        this.post(`/api/heads/${post.head.ID}`, post.headParameters).subscribe(
          (headParameters: APIHeadParameters) => {
            // do not update head from POST, wait for websocket...
            console.log(`POST head ${post.head.ID}: OK`, headParameters);
          }
          // TODO: errors to console?
        );
      }
    );

    this.webSocket = webSocketService.connect<APIEvents>('/events').subscribe(
      (apiEvents: APIEvents) => {
        console.log("WebSocket APIEvents", apiEvents);

        this.loadHeads(apiEvents.Heads);
      },
      (error: WebSocketError) => {
        console.log("WebSocket Error", error);
      },
      () => {
        console.log("WebSocket Close");
      }
    );
  }

  private loadHeads(apiHeads: APIHeads) {
    for (let id in apiHeads) {
      if (this.heads[id]) {
        this.heads[id].load(apiHeads[id]);
      } else {
        this.heads[id] = new Head(this.postSubject, apiHeads[id]);
      }
    }

    console.log("Loaded heads", this.heads);
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
