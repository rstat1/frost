import { MatSnackBar } from '@angular/material';
import { Subscription } from 'rxjs/Subscription';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { APIService } from 'app/services/api/api.service';
import { ActionListService, PrimaryActionInfo, SubActionClickEvent } from 'app/services/action-list.service';

@Component({
	selector: 'app-services',
	templateUrl: './services.html',
	styleUrls: ['./services.css']
})
export class ServicesRootComponent implements OnInit, OnDestroy {
	private actionClicked: Subscription;
	private subActionClicked: Subscription;
	private getServicesSub: Subscription;

	constructor(private actions: ActionListService, private route: ActivatedRoute,
		private router: Router, private api: APIService, private snackBar: MatSnackBar) {}

	ngOnInit() {
		console.warn("services root appears...")
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New Service", "add",
			"Add a new managed service"));
		this.actions.ClearSelectedItem();
		this.actions.SetSubItems([
			{IconName: "edit", Description: "Edit"},
			{IconName: "list", Description: "Logs"},
			{IconName: "delete", Description: "Delete"},
		])
		this.subActionClicked = this.actions.SubActionClicked.subscribe(action => {
			this.SubActionClicked(action);
		});
		this.getServicesSub = this.api.GetServices(true).subscribe(resp => {
			this.actions.SetActionList(JSON.parse(resp.response));
		});
		this.actionClicked = this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick = "New Service") {
				this.router.navigate(["new"], {
					relativeTo: this.route,
					skipLocationChange: true,
				});
			}
		})
	}
	ngOnDestroy(): void {
		console.warn("services root goes away...")
		this.actionClicked.unsubscribe();
		this.getServicesSub.unsubscribe();
		this.subActionClicked.unsubscribe();
	}
	private SubActionClicked(action: SubActionClickEvent) {
		if (action.ContextInfo != "watchdog") {
			switch(action.SubActionName) {
				case "Delete":
					this.api.DeleteService(action.ContextInfo).subscribe(resp => {
						this.actions.SetActionList(JSON.parse(resp.response));
					});
				break;
			}
		} else {
			this.snackBar.open("Can't delete the service you're using >_<", "", {
				duration: 3000, panelClass: "proper-colors", horizontalPosition: 'center',
				verticalPosition: 'top',
			});
		}
		console.log(action);
	}
	private handleActionClick(action: string) {
		if (action == "New Service") {
			this.router.navigate(["new"], {
				relativeTo: this.route
			});
		}
	}
}