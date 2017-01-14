import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import { WebSocketService, WebSocketError } from 'lib/websocket';

import * as _ from 'lodash';
import { Observable } from 'rxjs/Observable';
import { Observer } from 'rxjs/Observer';
import { Subscription } from 'rxjs/Subscription';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/map';

import { API, APIEvents, APIHeads, APIGroups, APIParameters, APIHeadParameters } from './api';
import { Post, Parameters, Head, Group } from './head';

@Injectable()
export class HeadService {
  private webSocket: Subscription;
  postSubject = new Subject<Post>();
  heads: Map<string, Head>;
  groups: Map<string, Group>;
  active: Head = null;

  listHeads(sort?: (Head) => any, filter?: (Head) => boolean): Head[] {
    let heads = Object.keys(this.heads).map(key => this.heads[key]);

    if (filter)
      heads = _.filter(heads, filter);

    if (sort)
      heads = _.sortBy(heads, sort);

    return heads;
  }
  listGroups(sort?: (Group) => any, filter?: (Groups) => boolean): Group[] {
    let groups = Array.from(this.groups.values());

    if (filter)
      groups = _.filter(groups, filter);

    if (sort)
      groups = _.sortBy(groups, sort);

    return groups;
  }

  byAddress(): Head[] {
    return this.listHeads(head => [head.Config.Universe, head.Config.Address]);
  }
  byID(): Head[] {
    return this.listHeads(head => head.ID);
  }

  ofIntensity(): Head[] {
    return this.listHeads(head => head.ID, head => !!head.Intensity);
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
    this.groups = new Map<string, Group>();

    this.get('/api/').subscribe(
      api => {
        this.loadHeads(api.Heads);
        this.loadGroups(api.Groups);
      }
    );

    this.postSubject.subscribe(
      post => {
        if (post.headID) {
          console.log(`POST head ${post.headID}...`, post.headParameters);

          this.post(`/api/heads/${post.headID}`, post.headParameters).subscribe(
            (headParameters: APIHeadParameters) => {
              // do not update head from POST, wait for websocket...
              console.log(`POST head ${post.headID}: OK`, headParameters);
            }
            // TODO: errors to console?
          );
        }

        if (post.groupID) {
          console.log(`POST group ${post.groupID}...`, post.groupParameters);

          this.post(`/api/groups/${post.groupID}`, post.groupParameters).subscribe(
            (groupParameters: APIParameters) => {
              // do not update head from POST, wait for websocket...
              console.log(`POST group ${post.groupID}: OK`, groupParameters);
            }
            // TODO: errors to console?
          );
        }
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

  private loadGroups(apiGroups: APIGroups) {
    for (let id in apiGroups) {
      let group: Group;

      if (group = this.groups.get(id)) {
        group.load(apiGroups[id]);
      } else {
        this.groups.set(id, new Group(this.postSubject, apiGroups[id]));
      }
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
