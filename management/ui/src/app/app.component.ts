import { Component } from '@angular/core';

@Component({
	selector: 'app-root',
	templateUrl: './app.component.html',
	styleUrls: ['./app.component.css']
})
export class AppComponent {
	title = 'app';
	constructor() {
		if (window.location.port == "4200") {
			document.title = "Watchdog-dev";
		} else if (window.location.hostname.includes("dev-m")) {
			document.title = "Watchdog-test";
		}
	}
}
