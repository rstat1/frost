import { Subscription } from 'rxjs';
import { MatSnackBar } from '@angular/material';
import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { APIService } from 'app/services/api/api.service';
import { ActionListService } from 'app/services/action-list/action-list.service';

@Component({
	selector: 'app-user-list',
	templateUrl: './user-list.component.html',
	styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit, OnDestroy {
	private getUserListSub: Subscription;
	private subActionClicked: Subscription;

	constructor(private actions: ActionListService, private api: APIService, private router: Router,
		private route: ActivatedRoute, private snackBar: MatSnackBar) { }

	ngOnInit() {
		this.actions.SetImageType(true);
		this.actions.SetActionList(null);
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
				this.router.navigate(["edit", action.ContextInfo], { relativeTo: this.route });
			}
		});
	}
	ngOnDestroy(): void {
		this.getUserListSub.unsubscribe();
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
