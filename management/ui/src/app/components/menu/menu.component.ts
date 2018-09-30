import { Component, OnInit } from '@angular/core';

import { MenuItem } from 'app/menu/menu-common';
import { MenuService } from 'app/services/menu.service';
import { AuthService } from 'app/services/auth/auth.service';

@Component({
	selector: 'app-menu',
	templateUrl: './menu.html',
	styleUrls: ['./menu.css']
})
export class MenuComponent implements OnInit {
	public isVisible: boolean = false;
	public categories: string[] = new Array();
	public menuItems: Map<string, Array<MenuItem>> = new Map();
	public scrollbarOptions = {
		scrollInertia: 0,
		theme: 'dark',
		scrollbarPosition: "inside",
		alwaysShowScrollbar: 0,
		autoHideScrollbar: true,
	};

	constructor(private menu: MenuService, private auth: AuthService) {}
	ngOnInit() {
		this.menu.GetMenuItems().forEach(item => {
			if (item.Icon == "") { item.Icon = "watchdog"; }
			if (this.menuItems.has(item.Category) == false) {
				if (item.Category == "Root" && this.auth.UserIsRoot) {
					this.menuItems.set(item.Category, [ item ]);
					this.categories.push(item.Category);
				}
				else if (item.Category != "Root") {
					this.menuItems.set(item.Category, [ item ]);
					this.categories.push(item.Category);
				}
			} else {
				let cat: MenuItem[] = this.menuItems.get(item.Category);
				if (item.RequiresRoot && this.auth.UserIsRoot) { cat.push(item); }
				else if (item.RequiresRoot == false || item.RequiresRoot == undefined) { cat.push(item); }
				this.menuItems.set(item.Category, cat);
			}
		});
	}
	public showMenu() {
		this.isVisible = !this.isVisible;
	}
	public getCategoryItems(category: string): MenuItem[] {
		let currentPage: string = this.menu.GetMenuContextData("currentPage");
		let items: MenuItem[] = this.menuItems.get(category);
		let result: MenuItem[] = new Array();
		items.forEach(item => {
			if (item.Context == currentPage || item.Context == null) {
				result.push(item);
			}
		});
		if (result.length > 0) { return result; }
		else { return null; }
	}
	public doSomethingWithClick(clickedItemTitle: string) {
		this.isVisible = false;
		this.menu.HandleMouseEvent(clickedItemTitle);
	}
}