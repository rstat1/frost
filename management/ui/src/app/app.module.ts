import { NgModule } from '@angular/core';
import { Routes, RouterModule } from "@angular/router";
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HTTP_INTERCEPTORS, HttpClientModule, HttpClient, HttpHandler } from '@angular/common/http';

import { AppComponent } from './app.component';
import { MenuItem } from "app/menu/menu-common";
import { APIService } from './services/api/api.service';
import { MenuService } from "app/services/menu.service";
import { AuthComponent } from 'app/components/auth/auth';
import { ManagerModule } from './manager/manager.module';
import { ConfigService } from './services/config.service';
import { AuthService } from 'app/services/auth/auth.service';
import { MenuModule, MenuItems } from 'app/menu/menu.module';
import { FirstRunComponent } from './manager/first-run/first-run';
import { AuthGuard, RootGuard } from './services/auth/auth.guard';
import { MenuComponent } from "app/components/menu/menu.component";
import { AuthTokenInjector } from './services/api/AuthTokenInjector';

const routes: Routes = [
	{path: 'auth', component: AuthComponent, pathMatch: "full"},
	{path: 'first-run', component: FirstRunComponent, pathMatch: "full"},
	{path: 'manage',  loadChildren: "app/manager/manager.module#ManagerModule", pathMatch: "full", canLoad: [AuthGuard]},
	{path: '**', redirectTo: "manage", pathMatch: 'full'}
];
const menuItems = { Items: [
	{ ItemTitle: "Users", ItemSubtext: "Return to Home page", Icon:"user", ActionName: "users", Category: "Config" },
	{ ItemTitle: "Updates", ItemSubtext: "Return to Home page", Icon:"update", ActionName: "update", Category: "Config" },
	{ ItemTitle: "Services", ItemSubtext: "Return to Home page", Icon:"services", ActionName: "services", Category: "Config" },
] };

@NgModule({
	declarations: [
		AppComponent,
		AuthComponent,
		FirstRunComponent,
	],
	imports: [
		BrowserModule,
		ManagerModule,
		HttpClientModule,
		BrowserAnimationsModule,
		MenuModule.forRoot(menuItems),
		RouterModule.forRoot(routes, {enableTracing: false}),
	],
	providers: [AuthService, APIService, AuthGuard, RootGuard, MenuService, ConfigService,
				{provide: HTTP_INTERCEPTORS, multi: true, useClass: AuthTokenInjector},
	],
	bootstrap: [AppComponent]
})
export class AppModule {
	constructor(private config: ConfigService) {}
}
