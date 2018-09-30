import { Router, ActivatedRoute } from '@angular/router';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnInit, Renderer2, OnDestroy } from '@angular/core';
import { MatTableDataSource, MatSnackBar } from '@angular/material';

import { APIService } from 'app/services/api/api.service';
import { AuthRequest, NewUser, ServiceAccess } from 'app/services/api/api-common';
import { ActionListService, PrimaryActionInfo } from 'app/services/action-list.service';
import { Subscription } from 'rxjs';

@Component({
	selector: 'app-users-root',
	templateUrl: './users-root.html',
	styleUrls: ['./users-root.css']
})
export class UsersRootComponent implements OnInit, OnDestroy {
	// private status: string = "";
	// private m: AuthRequest = new AuthRequest("","");
	// private dataSource = new MatTableDataSource<string>();
	// private displayedColumns = ['name', 'CanAccess', 'HasRoot'];
	// private selection = new SelectionModel<string>(true, []);
	// private permissions: ServiceAccess[] = new Array();
	private actionClicked: Subscription;
	private getUserListSub: Subscription;
	private subActionClicked: Subscription;

	constructor(private actions: ActionListService, private route: ActivatedRoute,
		private router: Router, private api: APIService) {}

	ngOnInit() {
		this.actions.SetImageType(true);
		this.actions.SetActionList(null);
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New User", "add", "Add a new user"));
		this.actionClicked = this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick == "New User") {
				this.router.navigate(["new"], {
					relativeTo: this.route
				});
			}
		});
		this.actions.SetSubItems([
			{IconName: "edit", Description: "Edit"},
			{IconName: "delete", Description: "Delete"},
		]);
		this.getUserListSub = this.api.GetUserList().subscribe(resp => {
			this.actions.SetActionList(JSON.parse(resp.response));
		});
		this.subActionClicked = this.actions.SubActionClicked.subscribe(action => {
			if (action.SubActionName == "delete") {
				this.api.DeleteUser(action.ContextInfo);
			}
		});
	}
	ngOnDestroy() {
		this.actionClicked.unsubscribe();
		this.getUserListSub.unsubscribe();
		this.subActionClicked.unsubscribe();
	}
}