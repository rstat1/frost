import { Component, OnInit, Input } from '@angular/core';

import { MenuItem } from 'app/menu/menu.module';
import { MenuService } from 'app/services/menu.service';
import { AuthService } from 'app/services/auth/auth.service';

@Component({
	selector: 'app-menu',
	templateUrl: './menu.html',
	styleUrls: ['./menu.css']
})
export class MenuComponent implements OnInit {
	@Input() public navType: string;
	@Input() public menuType: string = "app";

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

	constructor(public menu: MenuService, private auth: AuthService) {}
	ngOnInit() {
		this.categories = this.menu.GetCategoryList();
	}
	public showMenu() {
		this.isVisible = !this.isVisible;
	}
	public getCategoryItems(category: string): MenuItem[] {
		return this.menu.GetCategoryItems(category, this.menuType);
	}
	public doSomethingWithClick(clickedItemTitle: string) {
		if (clickedItemTitle == "home") {
			clickedItemTitle = "projects";
		}
		this.isVisible = false;
		this.menu.SetMenuContext(clickedItemTitle, "");
		this.menu.HandleMouseEvent(clickedItemTitle);
	}
}