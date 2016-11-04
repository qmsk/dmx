import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';
import { FormsModule }   from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';


import { AppComponent } from './app.component'
import { DashboardComponent } from './dashboard.component'
import { HeadsComponent } from './heads.component'

const routes: Routes = [
  { path: 'heads', component: HeadsComponent },
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
    HeadsComponent,
  ],
  bootstrap: [ AppComponent ],
})
export class AppModule {

}
