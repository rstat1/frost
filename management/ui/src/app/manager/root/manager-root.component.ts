import { Subscription } from 'rxjs';
import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { MenuService } from 'app/services/menu.service';
import { PageInfoService } from 'app/services/page-info.service';
import { ConfigService } from 'app/services/config.service';
import { APIService } from 'app/services/api/api.service';

@Component({
	selector: 'app-manager-root',
	templateUrl: './manager-root.component.html',
	styleUrls: ['./manager-root.component.css']
})
export class ManagerRootComponent implements OnInit, OnDestroy {
	public menuType: string = "app";
	public menuCategory: string = "App";
	public pageTitle: string = "frostcloud";
	public pageLogo: string = "watchdog";
	public showNeedsVT: boolean = false;
	private menuItemClicked: Subscription;

	constructor(private menu: MenuService, private router: Router, private pageInfo: PageInfoService, private route: ActivatedRoute, private api: APIService) { }

	ngOnInit() {
		console.log("manager root onInit");
		this.menu.CategoryChanged.subscribe(newCategory => {
			this.menuCategory = newCategory;
		});
		this.pageInfo.PageTitle.subscribe(newTitle => {
			this.pageTitle = newTitle;
		});
		this.pageInfo.PageLogo.subscribe(newLogoURL => {
			this.pageLogo = newLogoURL;
		});
		this.menuItemClicked = this.menu.MenuItemClicked.subscribe(item => {
			switch (item) {
				case "newservice":
					this.menu.SetMenuContext("newservice", "");
					this.router.navigate(["new"], { relativeTo: this.route });
					break;
			}
		});

		this.api.GetAppState().subscribe(resp => {
			if (resp.response == "initialized-need-vt") {
				this.showNeedsVT = true;
			}
		})
	}
	goHome() {
		this.router.navigate(["home"]);
	}
	ngOnDestroy(): void {
		this.menuItemClicked.unsubscribe();
	}
	public getServiceIconURL(name: string): string {
		return ConfigService.GetAPIURLFor("icon/" + name);
	}
}
