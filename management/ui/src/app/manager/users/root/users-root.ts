import { Router, ActivatedRoute } from '@angular/router';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnInit, Renderer2 } from '@angular/core';
import { MatTableDataSource, MatSnackBar } from '@angular/material';

import { APIService } from 'app/services/api/api.service';
import { AuthRequest, NewUser, ServiceAccess } from 'app/services/api/api-common';
import { ActionListService, PrimaryActionInfo } from 'app/services/action-list.service';

@Component({
	selector: 'app-users-root',
	templateUrl: './users-root.html',
	styleUrls: ['./users-root.css']
})
export class UsersRootComponent implements OnInit {
	private status: string = "";
	private m: AuthRequest = new AuthRequest("","");
	private dataSource = new MatTableDataSource<string>();
	private displayedColumns = ['name', 'CanAccess', 'HasRoot'];
	private selection = new SelectionModel<string>(true, []);
	private permissions: ServiceAccess[] = new Array();

	constructor(private actions: ActionListService, private route: ActivatedRoute,
		private router: Router) {}
		
	ngOnInit() {
		this.actions.SetActionList(null);
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New User", "add", "Add a new user"));
		this.actions.PrimaryActionClicked.subscribe(onClick => {
			if (onClick == "New User") {
				this.router.navigate(["new"], {
					relativeTo: this.route
				});
			}
		})
	}
}