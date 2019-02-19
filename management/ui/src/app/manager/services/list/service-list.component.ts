import { MatSnackBar } from '@angular/material';
import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';
import { ActionListService } from 'app/services/action-list/action-list.service';
import { SubActionClickEvent, PrimaryActionInfo } from 'app/services/action-list/action-list-common';

@Component({
	selector: 'app-service-list',
	templateUrl: './service-list.component.html',
	styleUrls: ['./service-list.component.css']
})
export class ServiceListComponent implements OnInit {
	private servicesList: string[];

	constructor(private actions: ActionListService, private api: APIService, private menu: MenuService,
		private router: Router, private route: ActivatedRoute, private snackBar: MatSnackBar) { }

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
		this.actions.SubActionClicked.subscribe(action => {
			this.SubActionClicked(action);
		});
		this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick = "New Service") {
				this.router.navigate(["new"], {
					relativeTo: this.route,
				});
			}
		});
		this.menu.SetMenuContext("list", "");
		this.menu.SetMenuCategory("App");
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
