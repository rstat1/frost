import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { CanActivate, Router, ActivatedRouteSnapshot,
		 RouterStateSnapshot, NavigationExtras } from '@angular/router';

import * as jwt_decode from 'jwt-decode';
import { environment } from 'environments/environment';
import { ConfigService } from 'app/services/config.service';
import { APIService } from 'app/services/api/api.service';
import { APIResponse } from 'app/services/api/api-common';

import 'rxjs/add/operator/toPromise'
import { HttpClient } from '@angular/common/http';

class SavedAuthDetails {
	public username: string;
	public token: string;
}

@Injectable()
export class AuthService {
	public CurrentUser: string = "";
	public RedirectURL: string = "";
	public FailureReason: string = "";
	public UserIsRoot: boolean = false;
	public IsLoggedIn: boolean = false;
	public NoToken: boolean = false;
	public AllowRegister: boolean = true;
	public AuthRequestInvalid: boolean = false;

	public AuthSuccess: Observable<boolean>;
	public TokenValidation: Observable<boolean>;

	private authSuccess: Subject<boolean>;
	private tokenValidate: Subject<boolean>;
	private savedAuthDetails: SavedAuthDetails;

	constructor(private api: APIService, private router: Router,
		private config: ConfigService, private http: HttpClient) {
		this.authSuccess = new Subject<boolean>();
		this.tokenValidate = new Subject<boolean>();

		this.AuthSuccess = this.authSuccess.asObservable();
		this.TokenValidation = this.tokenValidate.asObservable();

		this.setSavedToken()
	}
	public async setSavedToken() {
		let token: string = ConfigService.GetAccessToken();
		if (token != "") {
			let resp: APIResponse = await this.api.ValidateToken().toPromise().catch(e => {
				return null;
			})
			let decoded = jwt_decode(token);
			if (resp != null && resp.status == "success") {
				this.IsLoggedIn = true;
				this.UserIsRoot = decoded.grp == "root";
				this.tokenValidate.next(true);
			} else {
				sessionStorage.clear();
				this.IsLoggedIn = false;
				this.tokenValidate.next(false);
			}
		} else {
			this.IsLoggedIn = false;
			this.NoToken = true;
			this.tokenValidate.next(false);
		}
	}
	public doAuthRequest(username: string, password: string, redirect: string, isNewUser: boolean) {
		if (ConfigService.GetAccessToken() == "") {
			window.location.replace(ConfigService.GetAuthorizeEndpoint());
		} else {
			this.router.navigate(["manage"])
		}
	}
	public GetToken(authCode: string) {
		this.api.GetAuthToken(authCode).subscribe(resp => {
			this.handleAPIResponse(false, "manage", resp)
		});
	}
	private handleAPIResponse(isNewUser: boolean, redirectTo: string, resp: APIResponse) {
		let decoded = jwt_decode(resp.response);
		ConfigService.SetAccessToken(resp.response)

		let navigationExtras: NavigationExtras = {
			queryParamsHandling: 'preserve',
			preserveFragment: true
		};

		if (decoded.exp == null) {
			this.AuthRequestInvalid = true;
			this.CurrentUser = "";
			this.IsLoggedIn = false;
			this.FailureReason = "token not valid.";
		}
		else {
			this.UserIsRoot = decoded.grp == "root";
			this.CurrentUser = decoded.sub;
			sessionStorage.setItem("auth", JSON.stringify({username: this.CurrentUser, token: resp.response}));
			this.authSuccess.next(true);
			this.tokenValidate.next(true);
			this.NoToken = false
			this.IsLoggedIn = true;
		}

		// this.FailureReason = "";
		// this.IsLoggedIn = true;
		// this.AuthRequestInvalid = false;
		console.log(redirectTo)
		this.router.navigate([redirectTo]);
	}
	private handleAPIError(err: any) {
		this.AuthRequestInvalid = true;
		this.FailureReason = err.error.response
	}
}