import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/operator/map';

import { Head, ValueStream, HeadStream, StreamFunc, Channel, HeadIntensity, HeadColor } from './head';

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

  private streamHead(head: Head): StreamFunc {
    return (valueStream: ValueStream) => this.stream.next({head: head, valueStream: valueStream});
  }

  /*
   * Update head parameter values from GET/POST response.
   */
  private loadHead(head: Head, headData: Object) {
    let intensityData = headData['Intensity']; if (intensityData) {
      head.Intensity = new HeadIntensity(this.streamHead(head), intensityData);
    }
    let colorData = headData['Color']; if (colorData) {
      head.Color = new HeadColor(this.streamHead(head), colorData);
    }
  }

  private decodeHeads(headsData: Object[]): Head[] {
    let heads = headsData.map(headData => this.decodeHead(headData));
    heads.sort((a: Head, b: Head) => a.cmpHead(b));
    return heads;
  }
  private decodeHead(headData: Object): Head {
    let head = new Head;
    Object.assign(head, headData, {
      Channels: headData['Channels'].map(channelData => this.decodeChannel(channelData)),
    });
    this.loadHead(head, headData);

    return head;
  }
  private decodeChannel(channelData: Object): Channel {
    return Object.assign(new Channel, channelData);
  }

  load(): Observable<Head[]> {
    return this.http.get('/api/heads/')
      .map(response => this.heads = this.decodeHeads(response.json()))
    ;
  }

  setHeadChannel(head:Head, channel:Channel, params:Object): Observable<Channel> {
    return this.post(`/api/heads/${head.ID}/channels/${channel.ID}`, params)
      .map(channelParams => Object.assign(channel, channelParams));
  }
}
