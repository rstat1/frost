import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { Routes, RouterModule } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';

import { AppComponent } from 'app/app.component';
import { MenuService } from 'app/services/menu.service';
import { AuthComponent } from 'app/components/auth/auth';
import { APIService } from 'app/services/api/api.service';
import { ConfigService } from 'app/services/config.service';
import { MalihuScrollbarModule } from 'ngx-malihu-scrollbar';
import { AuthService } from 'app/services/auth/auth.service';
import { FirstRunComponent } from 'app/manager/first-run/first-run';
import { AuthGuard, RootGuard } from 'app/services/auth/auth.guard';
import { AuthTokenInjector } from 'app/services/api/AuthTokenInjector';

import { MenuModule } from 'app/menu/menu.module';
import { ManagerModule } from 'app/manager/manager.module';

const routes: Routes = [
	{path: 'auth', component: AuthComponent, pathMatch: "full"},
	{ path: 'services',  loadChildren: "app/manager/manager.module#ManagerModule", pathMatch: "full",
		canLoad: [AuthGuard]
	},
	{path: '', redirectTo: "/services", pathMatch: "full"}
];
const menuItems = { Items: [
	{ ItemTitle: "Users", ItemSubtext: "Edit user accounts", Icon:"user", Category:"App",
		ActionName: "users", RequiresRoot: false, MenuType: "app" },
	{ ItemTitle: "Telemetry", ItemSubtext: "View service telemetry data", Icon:"logs", Category:"App",
		ActionName: "services", RequiresRoot: false, MenuType: "app" },
	{ ItemTitle: "Updates", ItemSubtext: "Updates", Icon:"update", Category:"App",
		ActionName: "services", RequiresRoot: false, MenuType: "app" },

]};

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
		MalihuScrollbarModule.forRoot(),
		MenuModule.forRoot(menuItems),
		RouterModule.forRoot(routes),
	],
	providers: [AuthService, APIService, AuthGuard, RootGuard, MenuService, ConfigService,
		{provide: HTTP_INTERCEPTORS, multi: true, useClass: AuthTokenInjector}
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
