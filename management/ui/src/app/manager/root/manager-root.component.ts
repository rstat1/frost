import { Subscription } from 'rxjs';
import { Router } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { environment } from 'environments/environment';
import { MenuService } from 'app/services/menu.service';
import { PageInfoService } from 'app/services/page-info.service';

@Component({
	selector: 'app-manager-root',
	templateUrl: './manager-root.component.html',
	styleUrls: ['./manager-root.component.css']
})
export class ManagerRootComponent implements OnInit {
	public menuType: string = "app";
	public menuCategory: string = "App";
	public pageTitle: string = "frostcloud";
	public pageLogo: string = "watchdog";
	private getServicesSub: Subscription;

	constructor(private menu: MenuService, private router: Router, private pageInfo: PageInfoService) {}

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
	}
	goHome() {
		this.router.navigate(["home"]);
	}
	public getServiceIconURL(name: string): string {
		return environment.APIBaseURL + "/frost/icon/"+name;
	}
}
