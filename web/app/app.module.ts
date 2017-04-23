import { NgModule, ErrorHandler } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';
import { FormsModule }   from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';

import { StatusService } from './status.service'
import { AppErrorHandler } from './error-handler'

import { AppComponent } from './app.component'
import { ChannelsComponent } from './channels.component'
import { ColorComponent } from './color.component'
import { ColorControlsComponent } from './color-controls.component'
import { ColorPipe } from './color.pipe'
import { ControlComponent } from './control.component'
import { ControlsComponent } from './controls.component'
import { IntensityComponent } from './intensity.component'
import { PresetsComponent } from './presets.component'
import { PresetParametersComponent } from './preset-parameters.component'
import { MainComponent } from './main.component'
import { StatusComponent } from './status.component'
import { ValuePipe } from './value.pipe'

const routes: Routes = [
  { path: 'presets', component: PresetsComponent },
  { path: 'channels', component: ChannelsComponent },
  { path: 'intensity', component: IntensityComponent },
  { path: 'color', component: ColorComponent },
  { path: '', component: MainComponent },
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
    ChannelsComponent,
    ColorComponent,
    ColorControlsComponent,
    ColorPipe,
    ControlComponent,
    ControlsComponent,
    IntensityComponent,
    PresetsComponent,
    PresetParametersComponent,
    MainComponent,
    StatusComponent,
    ValuePipe,
  ],
  providers: [
    StatusService,
    {provide: ErrorHandler, useClass: AppErrorHandler},
  ],
  bootstrap: [ AppComponent ],
})
export class AppModule {

}
