import { catchError, map } from 'rxjs/operators';
import { of as observableOf, Observable } from 'rxjs';
import { Injectable } from "@angular/core";
import { CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot,
		 CanLoad, Route, CanActivateChild } from "@angular/router";
import { AuthService } from "app/services/auth/auth.service";

@Injectable()
export class AuthGuard implements CanActivate, CanLoad, CanActivateChild {
	constructor(private authService: AuthService, private router: Router) {}

	canActivateChild(childRoute: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
		let url = `/${childRoute}`;
		return this.checkLogin(url).pipe(
			map(status => {
				if (status == false) { this.router.navigate(['/auth']); }
				return status;
			}),
			catchError(e => observableOf(false)),);
		}
	canLoad(route: Route): boolean | Observable<boolean> | Promise<boolean> {
		let url = `/${route.path}`;
		return this.checkLogin(url).pipe(
			map(status => {
				if (status == false) { this.router.navigate(['/auth']); }
				return status;
			}),
			catchError(e => observableOf(false)),);
	}
	canActivate(route: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> {
		return this.checkLogin(routerState.url).pipe(
			map(status => {
				if (status == false) { this.router.navigate(['/auth']); }
				return status;
			}),
			catchError(e => observableOf(false)),);
	}
	checkLogin(redirectTo: string): Observable<boolean> {
		this.authService.RedirectURL = redirectTo;
		if (this.authService.IsLoggedIn) {
			return observableOf(true);
		} else if (this.authService.IsLoggedIn == false && this.authService.NoToken == false) {
			return this.authService.TokenValidation;
		} else {
			return observableOf(false);
		}
	}
}
@Injectable()
export class RootGuard implements CanActivate, CanLoad, CanActivateChild {
    constructor(private authService: AuthService, private router: Router) {}

    canActivateChild(childRoute: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
		return this.checkLogin(routerState.url).pipe(
			map(status => {
				return this.authService.UserIsRoot;
			}),
			catchError(e => observableOf(false)),);
    }
    canLoad(route: Route): boolean | Observable<boolean> | Promise<boolean> {
        let url = `/${route.path}`;
		return this.checkLogin(url).pipe(
			map(status => this.authService.UserIsRoot),
			catchError(e => observableOf(false)),);
    }
    canActivate(route: ActivatedRouteSnapshot, routerState: RouterStateSnapshot): boolean | Observable<boolean> | Promise<boolean> {
		return this.checkLogin(routerState.url).pipe(
			map(status => this.authService.UserIsRoot),
			catchError(e => observableOf(false)),);
    }
	checkLogin(redirectTo: string): Observable<boolean> {
		this.authService.RedirectURL = redirectTo;
		if (this.authService.IsLoggedIn) {
			return observableOf(true);
		} else if (this.authService.IsLoggedIn == false && this.authService.NoToken == false) {
			return this.authService.TokenValidation;
		} else {
			return observableOf(false);
		}
	}
}