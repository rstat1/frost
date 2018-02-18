import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpClient, HttpParams } from '@angular/common/http';

import { ConfigService } from "app/services/config.service";
import { APIResponse } from "app/services/api/api-common";

export class AuthRequest {
	public Username: string;
	public Password: string;
	constructor(username: string, password: string) {
		this.Username = username;
		this.Password = password;
	}
}

@Injectable()
export class APIService {
	constructor(private http: HttpClient) {}
	public GetServices(): Observable<APIResponse> {
		return this.GetRequest("services")
	}
	private GetRequest(url: string): Observable<APIResponse> {
		let apiURL: string = ConfigService.GetAPIURLFor(url)
		return this.http.get<APIResponse>(apiURL);
	}
	private PostRequest(url: string, body: string): Observable<APIResponse> {
		let apiURL: string = ConfigService.GetAPIURLFor(url)
		return this.http.post<APIResponse>(apiURL, body);
	}
}