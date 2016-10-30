import { Component, OnInit } from '@angular/core';

import { Head } from './head';
import { HeadService } from './head.service';


@Component({
  moduleId: module.id,
  selector: 'dmx-heads',
  templateUrl: 'heads.component.html',
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
        heads => this.heads = heads.sort((a: Head, b: Head) => a.cmpAddress(b)),
      )
    ;
  }
}
