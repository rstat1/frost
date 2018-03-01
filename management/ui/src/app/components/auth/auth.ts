import { Component, OnInit } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';

import { APIService } from 'app/services/api/api.service';
import { AuthService } from 'app/services/auth/auth.service';
import { ActivatedRoute } from '@angular/router';

@Component({
	selector: 'app-login',
	templateUrl: './auth.html',
	styleUrls: ['./auth.css']
})
export class AuthComponent {
	private getServicesSub: Subscription;
	private thisIsReallyDumb: Subscription;
	private superDuperDumb: Subscription;

	constructor(private api: APIService, private auth: AuthService, private route: ActivatedRoute) {
		let hasAuthCode: boolean = window.location.toString().includes("authcode");
		if (hasAuthCode == false) {
			this.getServicesSub = this.api.GetAppState().subscribe(resp => {
				if (resp.response == "initialized") {
					this.auth.doAuthRequest("", "", "", false);
				} else {
					window.location.replace(resp.response);
				}
			});
		} else if (hasAuthCode) {
			this.thisIsReallyDumb = this.route.url.subscribe(whyIsThisSoDumb => {
				if (whyIsThisSoDumb.length > 0) {
					if (whyIsThisSoDumb[0].path == "auth") {
						this.superDuperDumb = this.route.queryParams.subscribe(yUSoDumb => {
							console.log(yUSoDumb)
							this.auth.GetToken(yUSoDumb["code"]);
						});
					}
				}
			});
		}
	}
}