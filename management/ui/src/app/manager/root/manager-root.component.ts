import { Subscription } from 'rxjs';
import { MatSnackBar } from '@angular/material';
import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { APIService } from 'app/services/api/api.service';
import { ActionListService } from 'app/services/action-list/action-list.service';
import { PrimaryActionInfo, SubActionClickEvent } from 'app/services/action-list/action-list-common';

@Component({
	selector: 'app-manager-root',
	templateUrl: './manager-root.component.html',
	styleUrls: ['./manager-root.component.css']
})
export class ManagerRootComponent implements OnInit {
	private actionClicked: Subscription;
	private subActionClicked: Subscription;
	private getServicesSub: Subscription;
	private servicesList: string[];

	constructor(private actions: ActionListService, private api: APIService,
		private router: Router, private route: ActivatedRoute, private snackBar: MatSnackBar) {}

	ngOnInit() {
		console.log("manager root onInit");
		this.actions.SetImageType(false);
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New Service", "add",
			"Add a new managed service"));
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
		this.actionClicked = this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick = "New Service") {
				this.router.navigate(["new"], {
					relativeTo: this.route,
				});
			}
		});
	}
	SubActionClicked(action: SubActionClickEvent): any {
		switch(action.SubActionName) {
			case "Delete":
				if (action.ContextInfo != "watchdog") {
					this.api.DeleteService(action.ContextInfo).subscribe(resp => {
						this.servicesList.splice(this.servicesList.indexOf(action.ContextInfo), 1);
						this.actions.SetActionList(this.servicesList);
					});
				}
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
}
