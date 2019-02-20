import { Subscription } from 'rxjs';
import { Router } from '@angular/router';
import { Component, OnDestroy } from '@angular/core';

import { MenuService } from 'app/services/menu.service';

@Component({
	selector: 'app-root',
	templateUrl: './app.component.html',
	styleUrls: ['./app.component.css']
})
export class AppComponent implements OnDestroy {
	private menuItemClickedSub: Subscription;

	constructor(private router: Router, private menu: MenuService) {
		this.menuItemClickedSub = this.menu.MenuItemClicked.subscribe(itemClicked => {
			switch(itemClicked) {
				case "services":
					this.router.navigate(["services"]);
					break;
				case "vmconfig":
					// this.router.navigate([""])
					break;
				case "users":
					this.router.navigate(["users"]);
					break;
			}
		});
	}
	ngOnDestroy(): void {
		this.menuItemClickedSub.unsubscribe();
	}
}
