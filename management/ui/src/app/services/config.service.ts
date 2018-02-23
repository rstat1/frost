import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from "rxjs/Observable";
import { APIResponse } from "app/services/api/api-common";
import { environment } from '../../environments/environment';

@Injectable()
export class ConfigService {
	public static HOME_ROUTE_NAME: string = "/admin"; //Change when there's a home page route.
	private static API_VERSION_TAG: string = "frost";

	private static ACCESS_TOKEN: string = "";
	private static API_ENDPOINT: string = environment.APIBaseURL;
	private static AUTH_ENDPOINT: string = environment.APIBaseURL + "/trinity/";

	private instanceName : string;
	private instanceDescription : string;
	private registrationAllowed: boolean = null;

	constructor(private http: HttpClient) {}

	public static GetAPIURLFor(api: string, queryVars: string = ""): string {
		if (api == "ws") { return environment.WebsocketEndpoint + queryVars; }
		else { return this.API_ENDPOINT + "/" + this.API_VERSION_TAG + "/" + api; }
	}
	public static GetAuthURLFor(api: string) {
		return this.AUTH_ENDPOINT + api;
	}
	public static GetAuthorizeEndpoint(): string {
		return this.AUTH_ENDPOINT + "/authroize?sid=" + environment.ServiceID;
	}
	public static SetAPIEndpoint(endpoint: string) {
		this.API_ENDPOINT = endpoint;
	}
	public static GetAccessToken(): string
	{
		let savedAuthDetails = JSON.parse(sessionStorage.getItem("auth"))
		if (savedAuthDetails != null) { ConfigService.SetAccessToken(savedAuthDetails.token); }

		return this.ACCESS_TOKEN;
	}
	public static SetAccessToken(token: string) { this.ACCESS_TOKEN = token; }
	public get RegistrationAllowed(): boolean {
		console.log("hello")
		return this.registrationAllowed;
	}
	private GetConfigValue(valueName: string): Observable<APIResponse> {
		let apiURL: string = ConfigService.GetAPIURLFor("config/" + valueName)
		return this.http.get<APIResponse>(apiURL);
	}
	// public static getAuth0Callback() : string { return environment.Auth0Callback; }
	public static isProduction() : boolean { return environment.production; }
}
