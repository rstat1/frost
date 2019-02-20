import { Subscription } from 'rxjs';
import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { MenuService } from 'app/services/menu.service';

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

	constructor(private menu: MenuService, private router: Router, private route: ActivatedRoute) {}

	goHome() {
		this.router.navigate(["home"]);
	}
	ngOnInit() {
		this.menu.SetMenuContext("users", "");
		this.menu.CategoryChanged.subscribe(newCategory => {
			this.menuCategory = newCategory;
		});
		this.menu.MenuItemClicked.subscribe(item => {
			switch (item) {
				case "newuser":
					this.router.navigate(["new"], {relativeTo: this.route});
					break;
			}
		});
	}
	ngOnDestroy() {
		// this.actionClicked.unsubscribe();
		// this.getUserListSub.unsubscribe();
		// this.subActionClicked.unsubscribe();
	}
}