import { Subscription } from 'rxjs';
import { MatTableDataSource } from '@angular/material';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import "jquery";
import { Service } from 'app/services/api/api-common';
import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';
import { AuthService } from 'app/services/auth/auth.service';
import { ActionListService } from 'app/services/action-list.service';
import { PageInfoService } from 'app/services/page-info.service';

@Component({
	selector: 'app-manager-root',
	templateUrl: './manager-root.html',
	styleUrls: ['./manager-root.css']
})
export class ManagerRootComponent implements OnDestroy {
	public path: string;
	public pageName: string;
	public currentPath: string = "";
	public currentTitle: string = "";
	public showServiceList: boolean = false;
	public showSubActionArea: boolean = false;

	private lastAction: string = "";
	private lastService: string = "";
	private currentAction: string = "";
	private currentService: string = "";

	private pagePathSub: Subscription;
	private pageTitleSub: Subscription;
	private menuClickedSub: Subscription;
	private primaryActClickSub: Subscription;

	constructor(private menu: MenuService, private router: Router, private route: ActivatedRoute,
		private api: APIService, private auth: AuthService, private actions: ActionListService,
		private header: PageInfoService) {
		this.setHighlightsFromURL();
		this.menuClickedSub = this.menu.MenuItemClicked.subscribe(action => {
			this.menuItemClicked(action);
			if (this.showServiceList) {
				this.router.navigate([action], { relativeTo: this.route });
			}
		});
		this.primaryActClickSub = this.actions.PrimaryActionClicked.subscribe(() => {
			this.showSubActionArea = true;
		});
		this.actions.SubActionClicked.subscribe((action) => {
			if (action.SubActionName == "Edit" || action.SubActionName == "Logs") {
				this.showSubActionArea = true;
			}
		});
		this.pagePathSub = this.header.PagePath.subscribe(path => this.currentPath = path);
		this.pageTitleSub = this.header.PageTitle.subscribe(title => this.currentTitle = title);
	}
	ngOnDestroy() {
		this.menuClickedSub.unsubscribe();
	}
	private menuItemClicked(action: string) {
		this.showSubActionArea = false;
		if (this.lastAction == action || this.lastAction == undefined || !this.showServiceList) {
			this.showServiceList = !this.showServiceList;
		}
		this.setActionBackgroundColor(action);
		this.actions.ClearSelectedItem();
		if (!this.showServiceList) {
			$("#" + this.lastAction).css("background-color", "#1d1d1d");
		}
		this.router.navigate(["manage"]);
	}
	private setActionBackgroundColor(action: string) {
		this.currentAction = action;
		$("#" + action).css("background-color", "#272727");
		if (this.lastAction != "" && this.lastAction != action) {
			$("#" + this.lastAction).css("background-color", "#1d1d1d");
		}
		this.lastAction = action;
	}
	private setHighlightsFromURL() {
		//this.showSubActionArea = true;
		const urlBits: string[] = window.location.pathname.split("/");
		const specifiedService: string = window.location.search.replace("?service=", "");

		if (urlBits.length > 2) { this.showServiceList = true; }

		$("#" + urlBits[2]).ready(function() {
			$("#" + urlBits[2]).css("background-color", "#272727");
			this.showServiceList = true;
		});
		this.currentAction = this.lastAction = urlBits[2];
	}
	private navToAction(service: string) {
		this.router.navigate([this.currentAction], {
			relativeTo: this.route,
			queryParams: { service: service }
		});
	}
	public newService() {

	}
	public newUser() {
		$("#newuser").css("background-color", "#272727");
		this.router.navigate(["users/new"], {
			relativeTo: this.route
		});
	}
	public getCurrentAction(): string {
		return this.currentAction;
	}
}