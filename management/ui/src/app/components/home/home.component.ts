import { Router } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { MenuService } from 'app/services/menu.service';

@Component({
	selector: 'app-home',
	templateUrl: './home.component.html',
	styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

	constructor(public menu: MenuService, private router: Router) { }

	ngOnInit() {
		this.menu.SetMenuContext("home", "");
	}
	launch(what: string) {
		switch (what) {
			case "services":
				this.router.navigate(["services"]);
			break;
			case "users":
				this.router.navigate(["users"]);
			break;
		}
	}
	public getRowOrColumn(idx: number, type: string): number {
		if (type == "row") {
			if (idx == 0 || idx == 1) { return 1; }
			else { return 2; }
		} else if (type == "column") {
			if (idx == 0 || idx == 2) { return 1; }
			else { return 2; }
		}
	}
}
