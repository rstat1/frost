import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { Service } from 'app/services/api/api-common';
import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';

@Component({
  selector: 'app-manager-root',
  templateUrl: './manager-root.html',
  styleUrls: ['./manager-root.css']
})
export class ManagerRootComponent implements OnInit {
	public path: string;
	public pageName: string;
	public services: Service[] = [
		{"name": "player3", "filename":"", "api_prefix": "", "address": "", "managed": true},
		{"name": "gemini", "filename":"", "api_prefix": "", "address": "", "managed": true},
		{"name": "auth", "filename":"", "api_prefix": "", "address": "", "managed": true}
	];
	public showAddButton: boolean = false;
	public showServiceList: boolean = false;

	private lastAction: string = "";
	private currentAction: string = "";
	private knownActions: string[] = ["subprojectnew", "buildconfig", "todos", "projectnew"]

	private lastService: string = "";
	private currentService: string = "";

	ngOnInit() {}
	constructor(private menu: MenuService, private router: Router, private route: ActivatedRoute,
		private api: APIService) {
		this.setHighlightsFromURL();
		this.menu.MenuItemClicked.subscribe(action => {
			if (action == "services") { this.showAddButton = true; }
			else { this.showAddButton = false; }
			this.showServiceList = true;
			this.setActionBackgroundColor(action);
			this.router.navigate(["manage"])
		})
	}
	private setActionBackgroundColor(action: string) {
		this.currentAction = action;
		$("#"+action).css("background-color", "#272727");
		if (this.lastAction != "" && this.lastAction != action) {
			$("#"+this.lastAction).css("background-color", "#1d1d1d");
		}
		this.lastAction = action;
	}
	private setServiceBGColor(service: string) {
		this.currentService = service;
		$("#"+service).css("background-color", "#272727");
		if (this.lastService != "" && this.lastService != service) {
			$("#"+this.lastService).css("background-color", "#2e2e2e");
		}
		this.lastService = service;
	}
	private setHighlightsFromURL() {
		let urlBits: string[] = window.location.pathname.split("/");
		let specifiedService: string = window.location.search.replace("?service=", "");
		if (urlBits.length >= 3) {
			$("#"+urlBits[2]).ready(function() {
				$("#"+urlBits[2]).css("background-color", "#272727");
				$("#"+specifiedService).css("background-color", "#272727");
			});
			this.currentAction = this.lastAction = urlBits[2];
			this.currentService = this.lastService = specifiedService;
			this.showServiceList = true;
		}
	}
	private navToAction(service: string) {
		this.setServiceBGColor(service);
		this.router.navigate([this.currentAction], {
			relativeTo: this.route,
			queryParams: { service: service }
		});
	}
}