import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { MatSnackBar, MatDialog } from '@angular/material';

import { APIService } from 'app/services/api/api.service';
import { PageInfoService } from 'app/services/page-info.service';
import { Service, ServiceEdit, RouteAlias, AliasDeleteRequest } from 'app/services/api/api-common';
import { NewAliasDialogComponent } from 'app/manager/services/edit/new-alias-dialog/new-alias';

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
	public routeList: Map<string, Array<string>>;
	public apiRouteAliases: string[] = new Array();

	constructor(private header: PageInfoService, private route: ActivatedRoute,
				private api: APIService, private snackBar: MatSnackBar, private dialog: MatDialog) { }

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

				this.getAPIAliases();
			}
		});
	}
	public makeRouteList(resp: RouteAlias[]) {
		this.routeList = new Map();
		this.apiRouteAliases = new Array();
		if (resp.length > 0) {
			resp.forEach(item => {
				if (this.routeList.has(item.fullURL) == false) {
					this.routeList.set(item.fullURL, [item.apiRoute]);
					this.apiRouteAliases.push(item.fullURL);
				} else {
					let r: string[] = this.routeList.get(item.fullURL);
					r.push(item.apiRoute);
					this.routeList.set(item.fullURL, r);
				}
			});
		}
	}
	public getExtraRoutes(url: string): string[] {
		return this.routeList.get(url);
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
	public deleteAlias(routeToDel, fullURL: string) {
		let dar: AliasDeleteRequest = new AliasDeleteRequest;
		dar.route = routeToDel;
		dar.baseURL = fullURL;
		this.api.DeleteAlias(dar).subscribe(
			_ => {
				this.getAPIAliases();
				this.showResponse("Deleted alias successfully");
			},
			failed => this.showResponse("Failed: " + failed.error.response)
		);
	}
	public showNewAliasDialog() {
		let dialogRef = this.dialog.open(NewAliasDialogComponent, {
			width: "550px",
			data: {apiName: this.currentServiceAPIName},
		});
		dialogRef.afterClosed().subscribe(resp => {
			if (resp.status == "success") {
				this.showResponse("Added new route alias");
				this.getAPIAliases();
			}
		});
	}
	private getAPIAliases() {
		this.api.GetAPIAliases(this.currentServiceAPIName).subscribe(extras => {
			if (extras.status == "success") {
				 this.makeRouteList(JSON.parse(extras.response));
			}
		}, _ => {
			this.routeList = new Map();
			this.apiRouteAliases = new Array();
		});
	}
	private showResponse(message: string) {
		this.snackBar.open(message, "", {
			duration: 3000, panelClass: "proper-colors", horizontalPosition: 'right',
			verticalPosition: 'top',
		});
	}
}
