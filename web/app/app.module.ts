import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';
import { FormsModule }   from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';

import { AppComponent } from './app.component'
import { DashboardComponent } from './dashboard.component'
import { ChannelsComponent } from './channels.component'
import { ControlComponent } from './control.component'
import { HeadsComponent } from './heads.component'
import { HeadComponent } from './head.component'
import { IntensityComponent } from './intensity.component'

const routes: Routes = [
  { path: 'channels', component: ChannelsComponent },
  { path: 'intensity', component: IntensityComponent },
  { path: '', component: DashboardComponent },
];

@NgModule({
  imports: [
    BrowserModule,
    HttpModule,
    FormsModule,
    RouterModule.forRoot(routes, { useHash: true }),
  ],
  declarations: [
    AppComponent,
    DashboardComponent,
    ChannelsComponent,
    ControlComponent,
    HeadsComponent,
    HeadComponent,
    IntensityComponent,
  ],
  bootstrap: [ AppComponent ],
})
export class AppModule {

}
