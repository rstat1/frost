import { Subscription } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { APIService } from 'app/services/api/api.service';
import { AuthService } from 'app/services/auth/auth.service';
import { PageInfoService } from 'app/services/page-info.service';

@Component({
	selector: 'app-login',
	templateUrl: './auth.html',
	styleUrls: ['./auth.css']
})
export class AuthComponent implements OnInit {
	ngOnInit(): void {
		console.log("do auth...");
	}
	private getServicesSub: Subscription;
	private thisIsReallyDumb: Subscription;
	private superDuperDumb: Subscription;

	constructor(private api: APIService, private auth: AuthService, private route: ActivatedRoute, private pageInfo: PageInfoService) {
		let hasAuthCode: boolean = window.location.toString().includes("authcode");
		if (hasAuthCode == false) {
			this.getServicesSub = this.api.GetAppState().subscribe(resp => {
				if (resp.response == "initialized" || resp.response == "initialized-need-vt") {
					this.auth.doAuthRequest("", "", "", false);/*  */
				} else {
					window.location.replace(resp.response);
				}
				console.log(resp);
			}, e => { console.log(e); });
		} else if (hasAuthCode) {
			this.thisIsReallyDumb = this.route.url.subscribe(whyIsThisSoDumb => {
				if (whyIsThisSoDumb.length > 0) {
					if (whyIsThisSoDumb[0].path == "auth") {
						this.superDuperDumb = this.route.queryParams.subscribe(yUSoDumb => {
							this.auth.GetToken(yUSoDumb["code"]);
						});
					}
				}
			});
		}
	}
}