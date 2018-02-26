import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpClient, HttpParams } from '@angular/common/http';

import { AuthRequest } from './api-common'
import { ConfigService } from './config-service'

export interface APIResponse {
	status: string;
	response: string;
}

@Injectable()
export class APIService {
	constructor(private http: HttpClient) {}
	public GetSupportedServices(): Observable<APIResponse> {
		return this.FrostGetRequest("services");
	}
	public GetServiceName(rid: string): Observable<APIResponse> {
		return this.GetRequest("service/fromrid?r="+rid);
	}
	public ValidateCreds(request: AuthRequest, requestID: string): Observable<APIResponse> {
		let details = JSON.stringify(request);
		return this.PostRequest("validate?r="+requestID, details);
	}
	private FrostGetRequest(api: string) : Observable<APIResponse> {
		let frostURL: string = ConfigService.GetFrostURLFor(api)
		return this.http.get<APIResponse>(frostURL);
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