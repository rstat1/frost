import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from "rxjs";
import { APIResponse } from "app/services/api/api-common";
import { environment } from '../../environments/environment';

@Injectable()
export class ConfigService {
	public static HOME_ROUTE_NAME: string = "/admin"; //Change when there's a home page route.
	private static API_VERSION_TAG: string = "frost";

	private static ACCESS_TOKEN: string = "";
	private static SERVICE_ID: string = "";
	private static API_ENDPOINT: string = "";//environment.APIBaseURL;
	private static AUTH_ENDPOINT: string = "";

	private instanceName: string;
	private instanceDescription: string;
	private registrationAllowed: boolean = null;

	constructor(private http: HttpClient) {
		ConfigService.SetBaseURL();
		this.http.get<APIResponse>(ConfigService.GetAPIURLFor("serviceid")).subscribe(resp => {
			ConfigService.SERVICE_ID = resp.response;
		});
	}

	public static GetAPIURLFor(api: string, queryVars: string = ""): string {
		ConfigService.SetBaseURL();
		if (api == "ws") { return environment.WebsocketEndpoint + queryVars; }
		else { return ConfigService.API_ENDPOINT + "/" + ConfigService.API_VERSION_TAG + "/" + api; }
	}
	public static GetAuthURLFor(api: string) {
		ConfigService.SetBaseURL();
		return ConfigService.AUTH_ENDPOINT + api;
	}
	public static GetAuthorizeEndpoint(): string {
		ConfigService.SetBaseURL();
		return ConfigService.AUTH_ENDPOINT + "authorize?sid=" + ConfigService.SERVICE_ID;
	}
	public static SetAPIEndpoint(endpoint: string) {
		ConfigService.API_ENDPOINT = endpoint;
	}
	public static GetAccessToken(): string {
		if (sessionStorage.getItem("auth") != "") {
			const savedAuthDetails = JSON.parse(sessionStorage.getItem("auth"));
			if (savedAuthDetails != null) { ConfigService.SetAccessToken(savedAuthDetails.token); }
		}

		return this.ACCESS_TOKEN;
	}
	public static SetAccessToken(token: string) { this.ACCESS_TOKEN = token; }
	public get RegistrationAllowed(): boolean {
		return this.registrationAllowed;
	}
	private static SetBaseURL() {
		var s: string = window.location.hostname;
		if (s.startsWith("192") == false && s != "localhost") {
			var domain: string = s.substr(s.indexOf("."));
			ConfigService.API_ENDPOINT = window.location.protocol + "//" + "api" + domain;
		} else {
			ConfigService.API_ENDPOINT = "http://api.frostdev.m";
		}
		ConfigService.AUTH_ENDPOINT = ConfigService.API_ENDPOINT + "/trinity/";
	}
	private GetConfigValue(valueName: string): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAPIURLFor("config/" + valueName);
		return this.http.get<APIResponse>(apiURL);
	}
	// public static getAuth0Callback() : string { return environment.Auth0Callback; }
	public static isProduction(): boolean { return environment.production; }
}
