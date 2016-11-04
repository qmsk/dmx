import { Component, OnInit } from '@angular/core';

import { Head } from './head';
import { HeadService } from './head.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-heads',
  templateUrl: 'heads.component.html',
  styleUrls: [ 'heads.component.css' ],
  providers: [
    HeadService,
  ],
})
export class HeadsComponent implements OnInit {
  heads: Head[];

  constructor (private headService: HeadService) { }

  ngOnInit(): void {
    this.headService.list()
      .subscribe(
        heads => this.heads = heads.sort((a: Head, b: Head) => a.cmpHead(b)),
      )
    ;
  }

  select(head :Head) {
    head.active = !head.active;
    console.log("Toggle head active", head)
  }
}
