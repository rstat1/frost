import { Injectable, state } from "@angular/core";
import { CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot,
		 CanLoad, Route, CanActivateChild } from "@angular/router";
import { Observable } from "rxjs/Observable";
import { Subject } from 'rxjs/Subject';

import 'rxjs/add/operator/map'
import 'rxjs/add/operator/catch'
import 'rxjs/add/observable/of'

import { AuthService } from "app/services/auth/auth.service";
import { environment } from "environments/environment";

@Injectable()
export class AuthGuard implements CanActivate, CanLoad, CanActivateChild {
	constructor(private authService: AuthService, private router: Router) {}

	canActivateChild(childRoute: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
		let url = `/${childRoute}`;
		return this.checkLogin(url)
			.map(status => {
				if (status == false) { this.router.navigate(['/user/login']); }
				return status;
			})
			.catch(e => { return Observable.of(false); });
		}
	canLoad(route: Route): boolean | Observable<boolean> | Promise<boolean> {
		let url = `/${route.path}`;
		console.log("check login for " + url);
		return this.checkLogin(url)
			.map(status => {
				if (status == false) { this.router.navigate(['/user/login']); }
				return status;
			})
			.catch(e => { return Observable.of(false); });
	}
	canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean | Observable<boolean> {
		return this.checkLogin(state.url)
			.map(status => {
				console.log(status)
				if (status == false) { this.router.navigate(['/user/login']) }
				return status;
			})
			.catch(e => { return Observable.of(false); });
	}
	checkLogin(redirectTo: string): Observable<boolean> {
		this.authService.RedirectURL = redirectTo;
		if (this.authService.IsLoggedIn) {
			return Observable.of(true);
		} else if (this.authService.IsLoggedIn == false && this.authService.NoToken == false){
			console.log("wait")
			return this.authService.TokenValidation;
		} else {
			return Observable.of(false);
		}
	}
}
@Injectable()
export class RootGuard implements CanActivate, CanLoad, CanActivateChild {
    constructor(private authService: AuthService, private router: Router) {}

    canActivateChild(childRoute: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
        console.log("check login for " + state.url)
		return this.checkLogin(state.url)
			.map(status => {
				return this.authService.UserIsRoot;
			})
			.catch(e => { return Observable.of(false); });
    }
    canLoad(route: Route): boolean | Observable<boolean> | Promise<boolean> {
        let url = `/${route.path}`;
		console.log("check root for " + url)
		return this.checkLogin(url)
			.map(status => {
				console.log(status)
				return this.authService.UserIsRoot;
			})
			.catch(e => { return Observable.of(false); });
    }
    canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
        console.log("check login for " + state.url)
		return this.checkLogin(state.url)
			.map(status => {
				return this.authService.UserIsRoot;
			})
			.catch(e => { return Observable.of(false); });
    }
	checkLogin(redirectTo: string): Observable<boolean> {
		console.log(redirectTo)
		this.authService.RedirectURL = redirectTo;
		if (this.authService.IsLoggedIn) {
			return Observable.of(true);
		} else if (this.authService.IsLoggedIn == false && this.authService.NoToken == false){
			console.log("wait")
			return this.authService.TokenValidation;
		} else {
			return Observable.of(false);
		}
	}
}