import { Injectable, state } from "@angular/core";
import { CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot,
		 CanLoad, Route, CanActivateChild } from "@angular/router";
import { Observable } from "rxjs/Observable";

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/of';

import { AuthService } from "app/services/auth/auth.service";
import { environment } from "environments/environment";

@Injectable()
export class AuthGuard implements CanActivate, CanLoad, CanActivateChild {
	constructor(private authService: AuthService, private router: Router) {}

	canActivateChild(childRoute: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
		let url = `/${childRoute}`;
		return this.checkLogin(url)
			.map(status => {
				if (status == false) { this.router.navigate(['/auth']); }
				return status;
			})
			.catch(e => Observable.of(false));
		}
	canLoad(route: Route): boolean | Observable<boolean> | Promise<boolean> {
		let url = `/${route.path}`;
		return this.checkLogin(url)
			.map(status => {
				if (status == false) { this.router.navigate(['/auth']); }
				return status;
			})
			.catch(e => Observable.of(false));
	}
	canActivate(route: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> {
		return this.checkLogin(routerState.url)
			.map(status => {
				if (status == false) { this.router.navigate(['/auth']); }
				return status;
			})
			.catch(e => Observable.of(false));
	}
	checkLogin(redirectTo: string): Observable<boolean> {
		this.authService.RedirectURL = redirectTo;
		if (this.authService.IsLoggedIn) {
			return Observable.of(true);
		} else if (this.authService.IsLoggedIn == false && this.authService.NoToken == false) {
			return this.authService.TokenValidation;
		} else {
			return Observable.of(false);
		}
	}
}
@Injectable()
export class RootGuard implements CanActivate, CanLoad, CanActivateChild {
    constructor(private authService: AuthService, private router: Router) {}

    canActivateChild(childRoute: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
		return this.checkLogin(routerState.url)
			.map(status => {
				return this.authService.UserIsRoot;
			})
			.catch(e => Observable.of(false));
    }
    canLoad(route: Route): boolean | Observable<boolean> | Promise<boolean> {
        let url = `/${route.path}`;
		return this.checkLogin(url)
			.map(status => this.authService.UserIsRoot)
			.catch(e => Observable.of(false));
    }
    canActivate(route: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
		return this.checkLogin(routerState.url)
			.map(status => this.authService.UserIsRoot)
			.catch(e => Observable.of(false));
    }
	checkLogin(redirectTo: string): Observable<boolean> {
		this.authService.RedirectURL = redirectTo;
		if (this.authService.IsLoggedIn) {
			return Observable.of(true);
		} else if (this.authService.IsLoggedIn == false && this.authService.NoToken == false) {
			return this.authService.TokenValidation;
		} else {
			return Observable.of(false);
		}
	}
}