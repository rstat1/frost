import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpClient, HttpParams } from '@angular/common/http';

import { ConfigService } from "app/services/config.service";
import { APIResponse, NewUser } from "app/services/api/api-common";

@Injectable()
export class APIService {
	constructor(private http: HttpClient) {}
	public GetServices(minimal: boolean): Observable<APIResponse> {
		if (minimal) { return this.GetRequest("services?type=minimal"); }
		else { return this.GetRequest("services"); }
	}
	public TEMP_GetToken(): Observable<APIResponse> {
		return this.TrinityGetRequest("generatetoken");
	}
	public SaveUser(details: NewUser): Observable<APIResponse> {
		return this.TrinityPostRequest("user/new", JSON.stringify(details))
	}
	private TrinityGetRequest(endpoint: string): Observable<APIResponse> {
		let apiURL: string = ConfigService.GetAuthURLFor(endpoint);
		return this.http.get<APIResponse>(apiURL);
	}
	private TrinityPostRequest(endpoint: string, body: string): Observable<APIResponse> {
		let apiURL: string = ConfigService.GetAuthURLFor(endpoint);
		return this.http.post<APIResponse>(apiURL, body);
	}
	private GetRequest(url: string): Observable<APIResponse> {
		let apiURL: string = ConfigService.GetAPIURLFor(url);
		return this.http.get<APIResponse>(apiURL);
	}
	private PostRequest(url: string, body: string): Observable<APIResponse> {
		let apiURL: string = ConfigService.GetAPIURLFor(url);
		return this.http.post<APIResponse>(apiURL, body);
	}
}