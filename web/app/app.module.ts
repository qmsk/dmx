import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component'
import { HeadsComponent } from './heads.component'

@NgModule({
  imports: [
    BrowserModule,
    HttpModule,
  ],
  declarations: [
    AppComponent,
    HeadsComponent,
  ],
  bootstrap: [ AppComponent ],
})
export class AppModule {

}
