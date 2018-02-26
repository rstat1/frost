import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Routes, RouterModule } from "@angular/router";
import { BrowserModule } from '@angular/platform-browser';
import { MatButtonModule, MatInputModule, MatDialogModule,  MatDialog,
	MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { LoginComponent } from './login/login.component';
import { APIService } from './login/api-service';
import { ConfigService } from './login/config-service';
import { ErrorPageComponent } from './error-page/error-page.component';
import { AppComponent } from './app.component';

const routes: Routes = [
	{path: "login", component: LoginComponent, pathMatch: 'full'},
	{path: "error/:id", component: ErrorPageComponent, pathMatch: 'full'},
	{path: '',  pathMatch: 'full', redirectTo: "/login",}
]

@NgModule({
	declarations: [
		AppComponent,
		LoginComponent,
		ErrorPageComponent
	],
	imports: [
		FormsModule,
		BrowserModule,
		MatInputModule,
		MatButtonModule,
		HttpClientModule,
		BrowserAnimationsModule,
		RouterModule.forRoot(routes, {enableTracing: false}),
	],
	providers: [APIService, ConfigService],
	bootstrap: [AppComponent]
})
export class AppModule { }
