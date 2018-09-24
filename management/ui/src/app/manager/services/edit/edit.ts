import { MatSnackBar } from '@angular/material';
import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { Service, ServiceEdit, RouteAlias } from 'app/services/api/api-common';
import { APIService } from 'app/services/api/api.service';
import { PageInfoService } from 'app/services/page-info.service';

@Component({
	selector: 'app-edit',
	templateUrl: './edit.html',
	styleUrls: ['./edit.css']
})
export class EditServiceComponent implements OnInit {
	public aliasURL: string = "";
	public aliasedRoute: string = "";
	public currentARURL: string = "";
	public currentServiceID: string = "";
	public currentServiceKey: string = "";
	public currentServiceName: string = "";
	public currentLocalAddress: string = "";
	public currentServiceAPIName: string = "";
	public currentServiceAddress: string = "";
	public isCurrentServiceManaged: boolean;
	public isCurrentServiceUpdatesHosted: boolean;

	public uiFiles: File;
	public serviceBin: File;
	public uiFilesName: string = "";
	public s: Service = new Service();
	public serviceToEdit: string = "";

	constructor(private header: PageInfoService, private route: ActivatedRoute,
				private api: APIService, private snackBar: MatSnackBar) { }

	ngOnInit() {
		this.serviceToEdit = this.route.snapshot.paramMap.get('name');
		this.header.SetPagePath(window.location.pathname + "/edit");
		this.currentServiceName = this.serviceToEdit;
		this.api.GetService(this.serviceToEdit).subscribe(resp => {
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
	public setFile(name: string, event: any) {
		if (name == "ui") {
			this.uiFiles = event.target.files[0];
			if (this.uiFiles.type != "application/zip") {
				this.uiFiles = null;
				this.uiFilesName = "";
				this.snackBar.open("That wasn't a zip file. >_>", "", {
					duration: 3000, panelClass: "proper-colors", horizontalPosition: 'center',
					verticalPosition: 'top',
				});
			} else {
				this.uiFilesName = this.uiFiles.name;
			}
		} else {
			this.serviceBin = event.target.files[0];
			this.s.filename = this.serviceBin.name;
		}
	}
	public upload() {
		let hasValue: boolean = false;
		let body: FormData = new FormData();
		if (this.uiFiles != null) {
			hasValue = true;
			body.append("uiblob", this.uiFiles, this.uiFiles.name);
		}
		if (this.serviceBin != null) {
			hasValue = true;
			body.append("service", this.serviceBin, this.serviceBin.name);
		}
		if (hasValue) {
			this.api.UpdateService(body, this.currentServiceName).subscribe(
				resp => {
					this.snackBar.open("Update successful", "", {
						duration: 3000, panelClass: "proper-colors", horizontalPosition: 'right',
						verticalPosition: 'top',
					});
				},
				err => {
					this.snackBar.open(`Update failed: ${err.error.response}`, "", {
						duration: 3000, panelClass: "proper-colors", horizontalPosition: 'right',
						verticalPosition: 'top',
					});
				}
			);
		}
	}
	public save(propertyName: string) {
		let propChange: ServiceEdit = new ServiceEdit();
		propChange.property = propertyName;
		propChange.name = this.serviceToEdit;
		switch (propertyName) {
			case "name":
				propChange.new = this.currentServiceName;
				break;
			case "apiName":
				propChange.new = this.currentServiceAPIName;
				break;
			case "redirect":
				propChange.new = this.currentARURL;
				break;
			case "localaddr":
				propChange.new = this.currentLocalAddress;
				break;
			case "managed":
				if (this.isCurrentServiceManaged) {
					propChange.new = "Enabled";
				} else {
					propChange.new = "Disabled";
				}
				break;
		}
		this.api.EditService(propChange).subscribe(resp => {
			if (resp.status == "success") {
				if (resp.response != "success") { this.currentServiceKey = resp.response; }
				this.showResponse("Edit Successful");
			}
		}, e => this.showResponse("Edit failed: " + e.error.response) );
	}
	public newRoute() {
		let routeAlias: RouteAlias = new RouteAlias();
		routeAlias.apiName = this.currentServiceAPIName;
		routeAlias.fullURL = this.aliasURL;
		routeAlias.apiRoute = this.aliasedRoute;
		this.api.NewRouteAlias(routeAlias).subscribe(
			resp => {
				if (resp.status == "success") { this.showResponse("Added new route alias"); }
			},
			e => this.showResponse("Failed adding new alias: " + e.error.response)
		);
	}
	private showResponse(message: string) {
		this.snackBar.open(message, "", {
			duration: 3000, panelClass: "proper-colors", horizontalPosition: 'right',
			verticalPosition: 'top',
		});
	}
}
