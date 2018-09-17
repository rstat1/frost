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
	public GetService(name: string): Observable<APIResponse> {
		return this.GetRequest("service/get?name=" + name);
	}
	public GetAuthToken(code: string): Observable<APIResponse> {
		return this.GetRequest("auth/token?code=" + code);
	}
	public GetAppState(): Observable<APIResponse> {
		return this.GetRequest("status");
	}
	public GetServiceID(): Observable<APIResponse> {
		return this.GetRequest("serviceid");
	}
	public GetUserList(): Observable<APIResponse> {
		return this.TrinityGetRequest("user/list");
	}
	public InitWatchdog(): Observable<APIResponse> {
		return this.GetRequest("init");
	}
	public ValidateToken(): Observable<APIResponse> {
		return this.TrinityGetRequest("");
	}
	public SaveUser(details: NewUser): Observable<APIResponse> {
		return this.TrinityPostRequest("user/new", JSON.stringify(details));
	}
	public DeleteUser(username: string): Observable<APIResponse> {
		return this.TrinityDeleteRequest("user/delete?name=" + name);
	}
	public NewService(details: FormData): Observable<APIResponse> {
		return this.PostFormRequest("service/new", details);
	}
	public DeleteService(name: string): Observable<APIResponse> {
		return this.DeleteRequest("service/delete?name=" + name);
	}
	private TrinityGetRequest(endpoint: string): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAuthURLFor(endpoint);
		return this.http.get<APIResponse>(apiURL);
	}
	private TrinityPostRequest(endpoint: string, body: string): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAuthURLFor(endpoint);
		return this.http.post<APIResponse>(apiURL, body);
	}
	private TrinityDeleteRequest(endpoint: string): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAuthURLFor(endpoint);
		return this.http.delete<APIResponse>(apiURL);
	}
	private GetRequest(url: string): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAPIURLFor(url);
		return this.http.get<APIResponse>(apiURL);
	}
	private DeleteRequest(url: string): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAPIURLFor(url);
		return this.http.delete<APIResponse>(apiURL);
	}
	private PostRequest(url: string, body: string): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAPIURLFor(url);
		return this.http.post<APIResponse>(apiURL, body);
	}
	private PostFormRequest(url: string, body: FormData): Observable<APIResponse> {
		const apiURL: string = ConfigService.GetAPIURLFor(url);
		return this.http.post<APIResponse>(apiURL, body);
	}
}