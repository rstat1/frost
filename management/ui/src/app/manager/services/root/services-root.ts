import { Subscription } from 'rxjs/Subscription';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { ActionListService, PrimaryActionInfo } from 'app/services/action-list.service';
import { APIService } from 'app/services/api/api.service';

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
		private router: Router, private api: APIService) {
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New Service", "add",
			"Add a new managed service"));
	}

	ngOnInit() {
		this.actions.ClearSelectedItem();
		this.actions.SetSubItems([
			{IconName: "edit", Description: "Edit"},
			{IconName: "list", Description: "Logs"},
			{IconName: "delete", Description: "Delete"},
		])
	//	this.actionClicked = this.actions.PrimaryActionClicked.subscribe(this.handleActionClick);
		this.subActionClicked = this.actions.SubActionClicked.subscribe(this.handleSubActionClick);
		this.getServicesSub = this.api.GetServices(true).subscribe(resp => {
			this.actions.SetActionList(JSON.parse(resp.response));
		});
		this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick = "New Service") {
				this.router.navigate(["new"], {
					relativeTo: this.route
				});
			}
		})
	}
	ngOnDestroy(): void {
	//	this.actionClicked.unsubscribe();
		this.getServicesSub.unsubscribe();
	}
	private handleSubActionClick(action: string) {
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
