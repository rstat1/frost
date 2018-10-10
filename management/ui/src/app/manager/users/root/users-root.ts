import { Subscription } from 'rxjs';
import { MatSnackBar } from '@angular/material';
import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { APIService } from 'app/services/api/api.service';
import { ActionListService, PrimaryActionInfo } from 'app/services/action-list.service';

@Component({
	selector: 'app-users-root',
	templateUrl: './users-root.html',
	styleUrls: ['./users-root.css']
})
export class UsersRootComponent implements OnInit, OnDestroy {
	private actionClicked: Subscription;
	private getUserListSub: Subscription;
	private subActionClicked: Subscription;

	constructor(private actions: ActionListService, private route: ActivatedRoute,
		private router: Router, private api: APIService, private snackBar: MatSnackBar,) {}

	ngOnInit() {
		this.actions.SetImageType(true);
		this.actions.SetActionList(null);
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New User", "add", "Add a new user"));
		this.actionClicked = this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick == "New User") {
				this.router.navigate(["new"], {
					relativeTo: this.route,
					skipLocationChange: true,
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
			if (action.SubActionName == "Delete") {
				this.deleteUser(action.ContextInfo);
			} else if (action.SubActionName == "Edit") {
				this.router.navigate(["edit", action.ContextInfo], {
					relativeTo: this.route,
					skipLocationChange: true,
				});
			}
		});
	}
	ngOnDestroy() {
		this.actionClicked.unsubscribe();
		this.getUserListSub.unsubscribe();
		this.subActionClicked.unsubscribe();
	}
	private deleteUser(username: string) {
		this.api.DeleteUser(username).subscribe(_ => {
			this.snackBar.open("Deleted user successfully", "", {
				duration: 3000,
				panelClass: "proper-colors",
				horizontalPosition: 'right',
				verticalPosition: 'top'
			});
			this.getUserListSub.unsubscribe();
			this.getUserListSub = this.api.GetUserList().subscribe(r => {
				this.actions.SetActionList(JSON.parse(r.response));
			});
		}, error => {
			this.snackBar.open("Failed to delete user: " + error.error.response, "", {
				duration: 3000,
				panelClass: "proper-colors",
				horizontalPosition: 'right',
				verticalPosition: 'top'
			});
		});
	}
}