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
import { HomeComponent } from 'app/components/home/home.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { UsersModule } from './users/users.module';


const routes: Routes = [
	{ path: 'auth', component: AuthComponent, pathMatch: "full" },
	{
		path: 'services',
		loadChildren: "app/manager/manager.module#ManagerModule",
		pathMatch: "full",
		canActivate: [AuthGuard],
	},
	{
		path: 'users',
		loadChildren: "app/users/users.module#UsersModule",
		pathMatch: "full",
		canActivate: [AuthGuard],
	},
	{ path: 'home', component: HomeComponent, pathMatch: "full", canActivate: [AuthGuard] },
	{ path: '', redirectTo: "/home", pathMatch: "full" }
];
const menuItems = { Items: [
	//App Menu
	{ ItemTitle: "Users", ItemSubText: "Edit user accounts", Icon:"user", Category:"App",
		ActionName: "users", RequiresRoot: false, MenuType: "app", Context: "!users" },
	{ ItemTitle: "Default VM Config", ItemSubText: "Set default config for service VMs", Icon:"cloud", Category:"App",
		ActionName: "vmconfig", RequiresRoot: false, MenuType: "app" },
	{ ItemTitle: "Services", ItemSubText: "Edit or view service configuration and logs", Icon: "services", Category:"App",
		ActionName: "services", RequiresRoot: false, MenuType: "app", Context: "!services"},

	//New buttons
	{ ItemTitle: "New Service", ItemSubText:"Add a new managed service", Icon: "plus", Category: "App",
		ActionName: "newservice", RequiresRoot: false, MenuType: "app", Context: "services" },
	{ ItemTitle: "New User", ItemSubText:"Add a new managed service", Icon: "plus", Category: "App",
		ActionName: "newuser", RequiresRoot: false, MenuType: "app", Context: "users" },

	//Service menu
	{ ItemTitle: "Delete service", ItemSubText: "Delete this service", Icon: "delete", Category: "Service",
		MenuType: "app", ActionName: "deleteservice", RequiresRoot: false, Context: "service" },
	{ ItemTitle: "Service Logs", ItemSubText: "View service specific log data", Icon: "logs",
		MenuType: "app", ActionName: "logs", RequiresRoot: false, Category: "Service", Context: "service"  },
	{ ItemTitle: "VM Configuration", ItemSubText: "Edit configuration of the VM hosting this service", Icon: "cloud",
		MenuType: "app", ActionName: "editvmconfig", RequiresRoot: false, Category: "Service", Context: "service" },
	{ ItemTitle: "Restart Service", ItemSubText: "Stops and restarts this service", Icon: "restart",
		MenuType: "app", ActionName: "reboot", RequiresRoot: false, Category: "Service", Context: "service" },
	{ ItemTitle: "Configure Service", ItemSubText: "Allows specifying any service specifc configuration", Icon: "config",
		MenuType: "app", ActionName: "serviceconfig", RequiresRoot: false, Category: "Service", Context: "service" }
]};

@NgModule({
	declarations: [
		AppComponent,
		AuthComponent,
		HomeComponent,
		FirstRunComponent,
	],
	imports: [
		UsersModule,
		BrowserModule,
		ManagerModule,
		HttpClientModule,
		BrowserAnimationsModule,
		MenuModule.forRoot(menuItems),
		MalihuScrollbarModule.forRoot(),
		RouterModule.forRoot(routes, {enableTracing: false}),
	],
	providers: [AuthService, APIService, AuthGuard, RootGuard, MenuService, ConfigService,
		{provide: HTTP_INTERCEPTORS, multi: true, useClass: AuthTokenInjector}
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
