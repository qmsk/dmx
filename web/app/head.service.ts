import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/map';

import { Head, ValueStream, HeadStream, PostFunc, Channel, HeadIntensity, HeadColor } from './head';

@Injectable()
export class HeadService {
  stream = new Subject<HeadStream>();
  heads: Head[] = [];
  active: Head = null;

  select(head: Head) {
    console.log("Select head", head);
    this.active = head;
  }
  selected(head: Head): boolean {
    return this.active == head;
  }

  constructor(private http: Http) {
    this.stream.subscribe(
      headStream => {
        console.log(`Post head=${headStream.head.ID}`, headStream.valueStream);

        this.post(`/api/heads/${headStream.head.ID}`, headStream.valueStream).subscribe(
          headParams => {
            console.log(`Load head=${headStream.head.ID}`, headParams);

            this.loadHead(headStream.head, headParams);
          }
        );
      }
    );
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

  private postHeadFunc(head: Head): PostFunc {
    return (valueStream: ValueStream) => this.stream.next({head: head, valueStream: valueStream});
  }

  /*
   * Update head parameter values from GET/POST response.
   */
  private loadHead(head: Head, headData: Object) {
    let post = this.postHeadFunc(head);

    let channelsData = headData['Channels']; if (channelsData) {
      for (let channelID in channelsData) {
        let channel = head.channels[channelID]; if (channel) {
          head.channels[channelID].load(channelsData[channelID]);
        } else {
          head.channels[channelID] = new Channel(post, channelsData[channelID]);
        }
      }
    }
    let intensityData = headData['Intensity']; if (intensityData) {
      head.Intensity = new HeadIntensity(post, intensityData);
    }
    let colorData = headData['Color']; if (colorData) {
      head.Color = new HeadColor(post, colorData);
    }
  }
  private loadHeads(headsData: Object[]): Head[] {
    let heads = headsData.map(headData => {
      let head = new Head(headData);
      this.loadHead(head, headData);
      return head;
    });
    heads.sort((a: Head, b: Head) => a.cmpHead(b));
    return heads;
  }

  load(): Observable<Head[]> {
    return this.http.get('/api/heads/')
      .map(response => this.heads = this.loadHeads(response.json()))
    ;
  }

  setHeadChannel(head:Head, channel:Channel, params:Object): Observable<Channel> {
    return this.post(`/api/heads/${head.ID}/channels/${channel.ID}`, params)
      .map(channelParams => Object.assign(channel, channelParams));
  }
}
