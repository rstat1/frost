import { Subscription } from 'rxjs/Subscription';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import "jquery";
import { Service } from 'app/services/api/api-common';
import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'app-manager-root',
  templateUrl: './manager-root.html',
  styleUrls: ['./manager-root.css']
})
export class ManagerRootComponent implements OnInit, OnDestroy {
	public path: string;
	public pageName: string;
	public services: Service[] = [];
	public showAddButton: boolean = false;
	public showAddUser: boolean = false;
	public showServiceList: boolean = false;

	private lastAction: string = "";
	private currentAction: string = "";
	private knownActions: string[] = ["subprojectnew", "buildconfig", "todos", "projectnew"]
	private lastService: string = "";
	private currentService: string = "";
	private menuClickedSub: Subscription;
	private getServicesSub: Subscription

	ngOnInit() {

	}
	ngOnDestroy() {
		this.getServicesSub.unsubscribe();
		this.menuClickedSub.unsubscribe();
	}
	constructor(private menu: MenuService, private router: Router, private route: ActivatedRoute,
		private api: APIService, private auth: AuthService) {
		this.setHighlightsFromURL();
		this.auth.doAuthRequest("", "", "", false);
		this.menuClickedSub = this.menu.MenuItemClicked.subscribe(action => {
			this.showServiceList = true;
			this.setActionBackgroundColor(action);
			this.clearServiceBGColor();
			this.showActionButtons();
			$("#newuser").css("background-color", "#1d1d1d");
			$("#newservice").css("background-color", "#1d1d1d");
			if (this.currentAction == "services" || this.currentAction == "logs") {
				console.log(this.currentAction)
				this.getServicesSub = this.api.GetServices(false).subscribe(resp => {
					if (resp.status == "success") {
						this.services = JSON.parse(resp.response);
					}
				});
			} else { this.services = null; }
			this.router.navigate(["manage"])
		})
	}
	private setActionBackgroundColor(action: string) {
		this.currentAction = action;
		$("#"+action).css("background-color", "#272727");
		if (this.lastAction != "" && this.lastAction != action) {
			console.log(this.lastAction)
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
	private clearServiceBGColor() {
		if (this.lastService != "") {
			$("#"+this.lastService).css("background-color", "#2e2e2e");
		}
	}
	private setHighlightsFromURL() {
		let urlBits: string[] = window.location.pathname.split("/");
		let specifiedService: string = window.location.search.replace("?service=", "");

		if (urlBits.length > 2) { this.showServiceList = true; }

		$("#"+urlBits[2]).ready(function() {
			console.log(urlBits[2])
			if (urlBits.length >= 3 && urlBits[3] == "new") {
				// $("#users").css("background-color", "#272727");
				$("#newuser").css("background-color", "#272727");
			}
			$("#"+urlBits[2]).css("background-color", "#272727");
			this.showServiceList = true;
			// if (urlBits.length >= 3) {
			// 	$("#"+specifiedService).css("background-color", "#272727");
			// 	this.currentService = this.lastService = specifiedService;
			// }
		});
		this.currentAction = this.lastAction = urlBits[2];
		this.showActionButtons();
	}
	private navToAction(service: string) {
		this.setServiceBGColor(service);
		this.router.navigate([this.currentAction], {
			relativeTo: this.route,
			queryParams: { service: service }
		});
	}
	private showActionButtons() {
		if (this.currentAction == "services") { this.showAddButton = true; }
		else { this.showAddButton = false; }

		if (this.currentAction == "users" || this.currentAction == "newuser") {
			this.showAddUser = true;
			if (this.currentAction == "newuser") { }
		}
		else { this.showAddUser = false; }
	}
	public newService() {

	}
	public newUser() {
		$("#newuser").css("background-color", "#272727");
		this.router.navigate(["users/new"], {
			relativeTo: this.route
		});
	}
}