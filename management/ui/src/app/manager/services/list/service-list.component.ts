import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { MatSnackBar, MatDialog } from '@angular/material';

import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';
import { PageInfoService } from 'app/services/page-info.service';
import { ActionListService } from 'app/services/action-list/action-list.service';
import { SubActionClickEvent} from 'app/services/action-list/action-list-common';
import { DeleteServiceDialogComponent } from 'app/manager/services/edit/delete-service-dialog/delete-service-dialog';
import { Subscription } from 'rxjs';

@Component({
	selector: 'app-service-list',
	templateUrl: './service-list.component.html',
	styleUrls: ['./service-list.component.css']
})
export class ServiceListComponent implements OnInit, OnDestroy {
	private servicesList: string[];
	private subActionClicked: Subscription;

	constructor(private actions: ActionListService, private api: APIService, private menu: MenuService,
		private router: Router, private route: ActivatedRoute, private snackBar: MatSnackBar,
		private pageInfo: PageInfoService, private dialog: MatDialog) { }

	ngOnInit() {
		this.actions.SetImageType(false);
		this.actions.ClearSelectedItem();
		this.actions.SetSubItems([
			{IconName: "edit", Description: "Edit"},
			{IconName: "delete", Description: "Delete"},
			{IconName: "cached", Description: "Restart"},
		]);
		this.api.GetServices(true).subscribe(s => {
			this.servicesList = JSON.parse(s.response);
			this.actions.SetActionList(this.servicesList);
		});
		this.subActionClicked = this.actions.SubActionClicked.subscribe(action => {
			this.SubActionClicked(action);
		});
		this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick = "New Service") {
				this.router.navigate(["new"], {
					relativeTo: this.route,
				});
			}
		});
		this.menu.SetMenuContext("services", "");
		this.pageInfo.SetPageLogoAndTitle("watchdog","frostcloud");
		this.menu.SetMenuCategory("App");
	}
	ngOnDestroy(): void {
		this.subActionClicked.unsubscribe();
	}
	SubActionClicked(action: SubActionClickEvent): any {
		switch(action.SubActionName) {
			case "Delete":
				this.deleteService(action.ContextInfo);
				break;
			case "Edit":
				this.router.navigate([action.ContextInfo], { relativeTo: this.route });
				break;
			case "Restart":
				this.api.RestartService(action.ContextInfo).subscribe(resp => {
					if (resp.status == "success") {
						this.snackBar.open("Restart successful", "", {
							duration: 3000, panelClass: "proper-colors", horizontalPosition: 'center',
							verticalPosition: 'top',
						});
					}
				});
				break;
		}
	}
	private deleteService(serviceName: string) {
		console.log('delete service NAO!');
		this.dialog.open(DeleteServiceDialogComponent, {
			width:'500px',
			data: {project: serviceName},
		}).afterClosed().subscribe(res => {
			if (res == true) {
				this.servicesList.splice(this.servicesList.indexOf(serviceName), 1);
				this.actions.SetActionList(this.servicesList);
			}
		});
	}
}
