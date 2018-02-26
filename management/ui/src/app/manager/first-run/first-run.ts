import { Component, OnInit, OnDestroy } from '@angular/core';

import { APIService } from 'app/services/api/api.service';
import { ConfigService } from "app/services/config.service";

class InstanceInfo {
	public Password: string
	public ServiceID: string
	public ServiceKey: string
}

@Component({
	selector: 'app-first-run',
	templateUrl: './first-run.html',
	styleUrls: ['./first-run.css']
})
export class FirstRunComponent implements OnInit, OnDestroy {
	private key: string = "";
	private id: string = "";
	private rootPW: string = "";

	constructor(private api: APIService) {
		(<any>document.getElementsByClassName("background")[0]).style.backgroundImage = 'url("' +
			ConfigService.GetAPIURLFor("bg") + '")'
		this.api.InitWatchdog().subscribe(resp => {
			this.rootPW = resp.response;
		});
	}
}