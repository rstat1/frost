import { Subscription } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatSnackBar, MatDialog } from '@angular/material';

import { environment } from 'environments/environment';
import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';
import { ConfigService } from 'app/services/config.service';
import { PageInfoService } from 'app/services/page-info.service';
import { NewAliasDialogComponent } from 'app/manager/services/edit/new-alias-dialog/new-alias';
import { Service, ServiceEdit, RouteAlias, AliasDeleteRequest } from 'app/services/api/api-common';
import { DeleteServiceDialogComponent } from 'app/manager/services/edit/delete-service-dialog/delete-service-dialog';

@Component({
	selector: 'app-edit',
	templateUrl: './edit.html',
	styleUrls: ['./edit.css']
})
export class EditServiceComponent implements OnInit, OnDestroy {
	public aliasURL: string = "";
	public aliasedRoute: string = "";
	public currentARURL: string = "";
	public currentFileName: string = "";
	public currentServiceID: string = "";
	public currentServiceKey: string = "";
	public currentServiceName: string = "";
	public currentLocalAddress: string = "";
	public currentServiceAPIName: string = "";
	public currentServiceAddress: string = "";
	public enableVaultIntegration: boolean;
	public isCurrentServiceManaged: boolean;
	public isCurrentServiceUpdatesHosted: boolean;

	public icon: File;
	public uiFiles: File;
	public show: boolean;
	public success: boolean;
	public serviceBin: File;
	public iconName: string = "";
	public lastError: string = "";
	public uiFilesName: string = "";
	public s: Service = new Service();
	public serviceToEdit: string = "";
	public routeList: Map<string, Array<string>>;
	public apiRouteAliases: string[] = new Array();

	private menuItemClickedSub: Subscription;

	constructor(private header: PageInfoService, private route: ActivatedRoute, private menu: MenuService,
		private api: APIService, private snackBar: MatSnackBar, private dialog: MatDialog,
		private pageInfo: PageInfoService) { }

	ngOnInit() {
		this.serviceToEdit = this.route.snapshot.paramMap.get('name');
		this.header.SetPagePath(window.location.pathname + "/edit");
		this.currentServiceName = this.serviceToEdit;
		this.api.GetService(this.serviceToEdit).subscribe(resp => {
			if (resp.status == "success") {
				let service: Service = JSON.parse(resp.response);
				this.currentFileName = service.filename;
				this.currentARURL = service.RedirectURL;
				this.currentServiceID = service.ServiceID;
				this.currentLocalAddress = service.address;
				this.isCurrentServiceManaged = service.managed;
				this.currentServiceAPIName = service.api_prefix;
				this.isCurrentServiceUpdatesHosted = service.managedUpdates;
				this.enableVaultIntegration = service.needsVault;

				this.getAPIAliases();
			}
		});
		this.menu.SetMenuContext("service", "");
		this.menu.SetMenuCategory("Service");
		this.pageInfo.SetPageLogoAndTitle(this.serviceToEdit, this.serviceToEdit);
		this.menuItemClickedSub = this.menu.MenuItemClicked.subscribe(item => {
			switch (item) {
				case "reboot":
					if (this.serviceToEdit != "watchdog") {
						this.api.RestartService(this.serviceToEdit).subscribe(resp => {
							if (resp.status == "success") {
								this.showResponse("Restart successful.");
							}
						});
					}
					break;
				case "deleteservice":
					if (this.serviceToEdit != "watchdog") {
						this.deleteService();
					} else {
						this.showResponse("Cannot delete the watchdog service.");
					}
					break;
				case "vmconfig":
					break;
				case "logs":
					break;
				case "serviceconfig":
					break;
			}
		});
	}
	ngOnDestroy() {
		this.menuItemClickedSub.unsubscribe();
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
		} else if (name == "service") {
			this.serviceBin = event.target.files[0];
			this.s.filename = this.serviceBin.name;
		} else if (name == "icon") {
			this.icon = event.target.files[0];
			this.iconName = this.icon.name;
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
			this.show = false;
			this.api.UpdateService(body, this.currentServiceName).subscribe(
				resp => {
					this.snackBar.open("Update successful", "", {
						duration: 3000, panelClass: "proper-colors", horizontalPosition: 'right',
						verticalPosition: 'top',
					});
					this.show = true;
					this.success = true;
				},
				err => {
					this.snackBar.open(`Update failed: ${err.error.response}`, "", {
						duration: 3000, panelClass: "proper-colors", horizontalPosition: 'right',
						verticalPosition: 'top',
					});
					this.show = true;
					this.success = false;
					this.lastError = err.error.response;
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
			case "filename":
				propChange.new = this.currentFileName;
				break;
			case "icon":
				propChange.new = "icon";
				break;
			case "vault":
				if (this.enableVaultIntegration) {
					propChange.new = "Enabled";
				} else {
					propChange.new = "Disabled";
				}
				console.log(propChange)
				break;
		}
		if (propChange.new != "icon") {
			this.api.EditService(propChange).subscribe(resp => {
				if (resp.status == "success") {
					if (resp.response != "success") { this.currentServiceKey = resp.response; }
					this.showResponse("Edit Successful");
				}
			}, e => this.showResponse("Edit failed: " + e.error.response));
		} else {
			this.uploadIcon();
		}
	}
	public uploadIcon() {
		let body: FormData = new FormData();
		if (this.icon != null) {
			body.append("icon", this.icon, this.icon.name);
			this.api.UploadIcon(body, this.currentServiceName).subscribe(resp => {
				this.showResponse("Upload Successful");
			}, e => this.showResponse("Failed: " + e.error.response));
		} else {
			this.showResponse("Pick an icon first.");
		}
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
			data: { apiName: this.currentServiceAPIName },
		});
		dialogRef.afterClosed().subscribe(resp => {
			if (resp.status == "success") {
				this.showResponse("Added new route alias");
				this.getAPIAliases();
			}
		});
	}
	public getServiceIconURL(name: string): string {
		return ConfigService.GetAPIURLFor("icon/" + name);
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
	private deleteService() {
		this.dialog.open(DeleteServiceDialogComponent, {
			width: '500px',
			data: { project: this.serviceToEdit },
		});
	}
}
