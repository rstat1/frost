import { Observable } from 'rxjs';
import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';

import { ConfigService } from "app/services/config.service";
import { APIResponse, NewUser, ServiceEdit, RouteAlias, AliasDeleteRequest, PermissionChange, PasswordChange } from "app/services/api/api-common";

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
	public GetUserInfo(username: string): Observable<APIResponse> {
		return this.TrinityGetRequest("user?name=" +username);
	}
	public GetPermissionMap(username: string): Observable<APIResponse> {
		return this.TrinityGetRequest("permissions?user=" + username);
	}
	public GetAPIAliases(api: string): Observable<APIResponse> {
		return this.GetRequest("aliases/all?api="+api);
	}
	public ChangePermissionValue(change: PermissionChange): Observable<APIResponse> {
		return this.TrinityPostRequest("permissions/change", JSON.stringify(change));
	}
	public ChangePassword(change: PasswordChange): Observable<APIResponse> {
		return this.TrinityPostRequest("user/edit", JSON.stringify(change));
	}
	public ValidateToken(): Observable<APIResponse> {
		return this.TrinityGetRequest("");
	}
	public InitWatchdog(): Observable<APIResponse> {
		return this.GetRequest("init");
	}
	public DeleteUser(username: string): Observable<APIResponse> {
		return this.TrinityDeleteRequest("user/delete?name=" + username);
	}
	public DeleteService(name: string): Observable<APIResponse> {
		return this.DeleteRequest("service/delete?name=" + name);
	}
	public DeleteAlias(name: AliasDeleteRequest): Observable<APIResponse> {
		return this.PostRequest("aliases/delete", JSON.stringify(name));
	}
	public SaveUser(details: NewUser): Observable<APIResponse> {
		return this.TrinityPostRequest("user/new", JSON.stringify(details));
	}
	public EditService(propChange: ServiceEdit): Observable<APIResponse> {
		return this.PostRequest("service/edit", JSON.stringify(propChange));
	}
	public UpdateService(details: FormData, name: string): Observable<APIResponse> {
		return this.PostFormRequest("service/update?name="+name, details);
	}
	public NewRouteAlias(alias: RouteAlias): Observable<APIResponse> {
		return this.PostRequest("aliases/new", JSON.stringify(alias));
	}
	public NewService(details: FormData): Observable<APIResponse> {
		return this.PostFormRequest("service/new", details);
	}
	public UploadIcon(details: FormData, name: string): Observable<APIResponse> {
		return this.PostFormRequest("icon/new/"+name, details);
	}
	public RestartService(serviceName: string): Observable<APIResponse> {
		return this.GetRequest("service/restart/"+serviceName);
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