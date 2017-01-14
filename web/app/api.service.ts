import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import { WebSocketService, WebSocketError } from 'lib/websocket';

import * as _ from 'lodash';
import { Observable } from 'rxjs/Observable';
import { Observer } from 'rxjs/Observer';
import { Subscription } from 'rxjs/Subscription';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/map';

import { API, APIEvents, APIHeads, APIGroups, APIPresets } from './api';
import { Post, Head, Group, Preset } from './head';

@Injectable()
export class APIService {
  private webSocket: Subscription;
  postSubject = new Subject<Post>();

  // state
  heads: Map<string, Head>;
  groups: Map<string, Group>;
  presets: Map<string, Preset>;

  listHeads(sort?: (Head) => any, filter?: (Head) => boolean): Head[] {
    let heads = Array.from(this.heads.values());

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
  listPresets(sort?: (Preset) => any, filter?: (Preset) => boolean): Preset[] {
    let presets = Array.from(this.presets.values());

    if (filter)
      presets = _.filter(presets, filter);

    if (sort)
      presets = _.sortBy(presets, sort);

    return presets;
  }

  constructor(private http: Http, webSocketService: WebSocketService) {
    this.heads = new Map<string, Head>();
    this.groups = new Map<string, Group>();
    this.presets = new Map<string, Preset>();

    this.get('/api/').subscribe(
      api => {
        this.loadHeads(api.Heads);
        this.loadGroups(api.Groups);
        this.loadPresets(api.Presets);
      }
    );

    this.postSubject.subscribe(
      post => {
        console.log(`POST ${post.type} ${post.id}...`, post.parameters);

        this.post(`/api/${post.type}/${post.id}`, post.parameters).subscribe(
          (response) => {
            // do not update head from POST, wait for websocket...
            console.log(`POST ${post.type} ${post.id} OK`, response);
          }
          // TODO: errors to console?
        );
      }
    );

    this.webSocket = webSocketService.connect<APIEvents>('/events').subscribe(
      (apiEvents: APIEvents) => {
        console.log("WebSocket APIEvents", apiEvents);

        this.loadHeads(apiEvents.Heads);
        this.loadGroups(apiEvents.Groups);
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
      let head = this.heads.get(id)

      if (head) {
        head.load(apiHeads[id])
      } else {
        this.heads.set(id, new Head(this.postSubject, apiHeads[id]))
      }
    }

    console.log("Loaded heads", this.heads);
  }

  private loadGroups(apiGroups: APIGroups) {
    for (let id in apiGroups) {
      let group: Group;

      let heads = apiGroups[id].Heads.map((id) => this.heads.get(id));

      if (group = this.groups.get(id)) {
        group.load(apiGroups[id]);
      } else {
        this.groups.set(id, new Group(this.postSubject, apiGroups[id], heads));
      }
    }
  }

  private loadPresets(apiPresets: APIPresets) {
    for (let id in apiPresets) {
      let heads: Head[] = null
      let groups: Group[] = null

      if (apiPresets[id].Config.Heads) {
        heads = Object.keys(apiPresets[id].Config.Heads).map((id) => {
          return this.heads.get(id)
        })
      }

      if (apiPresets[id].Config.Groups) {
        groups = Object.keys(apiPresets[id].Config.Groups).map((id) => {
          return this.groups.get(id)
        })
      }

      this.presets.set(id, new Preset(this.postSubject, apiPresets[id], heads, groups));
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
