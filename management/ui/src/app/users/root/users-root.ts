import { Subscription } from 'rxjs';
import { Router } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';
import { PrimaryActionInfo } from 'app/services/action-list/action-list-common';
import { ActionListService } from 'app/services/action-list/action-list.service';

@Component({
	selector: 'app-users-root',
	templateUrl: './users-root.html',
	styleUrls: ['./users-root.css']
})
export class UsersRootComponent implements OnInit, OnDestroy {
	public menuType: string = "app";
	public menuCategory: string = "App";
	private actionClicked: Subscription;
	private getUserListSub: Subscription;
	private subActionClicked: Subscription;

	constructor(private menu: MenuService, private router: Router) {}

	goHome() {
		this.router.navigate(["home"]);
	}
	ngOnInit() {
		this.menu.SetMenuContext("users", "");
		this.menu.CategoryChanged.subscribe(newCategory => {
			this.menuCategory = newCategory;
		});
	}
	ngOnDestroy() {
		// this.actionClicked.unsubscribe();
		// this.getUserListSub.unsubscribe();
		// this.subActionClicked.unsubscribe();
	}
}