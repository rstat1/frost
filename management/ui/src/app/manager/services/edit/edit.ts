import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { Service } from 'app/services/api/api-common';
import { APIService } from 'app/services/api/api.service';
import { PageInfoService } from 'app/services/page-info.service';

@Component({
	selector: 'app-edit',
	templateUrl: './edit.html',
	styleUrls: ['./edit.css']
})
export class EditServiceComponent implements OnInit {
	public currentARURL: string = "";
	public currentServiceID: string = "";
	public currentServiceKey: string = "";
	public currentServiceName: string = "";
	public currentLocalAddress: string = "";
	public currentServiceAPIName: string = "";
	public currentServiceAddress: string = "";
	public isCurrentServiceManaged: boolean;
	public isCurrentServiceUpdatesHosted: boolean;

	public uiFilesName: string = "";
	public s: Service = new Service();

	constructor(private header: PageInfoService, private route: ActivatedRoute,
				private api: APIService) { }

	ngOnInit() {
		let serviceName = this.route.snapshot.paramMap.get('name');
		this.header.SetPagePath(window.location.pathname + "/edit");
		this.currentServiceName = serviceName;
		this.api.GetService(serviceName).subscribe(resp => {
			if (resp.status == "success") {
				let service: Service = JSON.parse(resp.response);
				this.currentARURL = service.RedirectURL;
				this.currentServiceID = service.ServiceID;
				this.currentLocalAddress = service.address;
				this.isCurrentServiceManaged = service.managed;
				this.currentServiceAPIName = service.api_prefix;
				this.isCurrentServiceUpdatesHosted = service.managedUpdates;
			}
		});
	}
	public setFile(name: string, event: any) {}
}
