import { Component } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { APIService } from './api-service'
import { AuthRequest } from './api-common'
import { ConfigService } from './config-service';

@Component({
	selector: 'app-login',
	templateUrl: './login.html',
	styleUrls: ['./login.css']
})
export class LoginComponent {
	public m: AuthRequest = new AuthRequest("", "");
	public invalid: boolean = false;
	public errorMessage: string = "";
	public serviceName: string = "(invalid)";
	public knownApps: string[] = []; //["player3", "gemini", "watchdog"]

	private requestID: string = "";

	constructor(private router: Router, private api: APIService) {
		this.requestID = window.location.search.replace("?r=", "");
		this.api.GetServiceName(this.requestID).subscribe(
			resp => this.serviceName = resp.response,
			error => {
				console.log(error);
				this.router.navigate(["/error", "1"]);
			}
		);
		this.api.GetSupportedServices().subscribe(resp => {
			this.knownApps = JSON.parse(resp.response);
		});
		(<any>document.getElementsByClassName("background")[0]).style.backgroundImage = 'url("' + ConfigService.GetFrostURLFor("bg") + '")'
	}
	public login() {
		let errorCode: string = "1";
		this.api.ValidateCreds(this.m, this.requestID).subscribe(
			resp => {
				window.location.replace(resp.response)
			},
			error => {
				if (error.error.response == "not authorized") { errorCode = "2"; }
				else if (error.error.response == "incorrect credentials") { errorCode = "3"; }
				this.router.navigate(["/error", errorCode]);
			}
		)
	}
}
