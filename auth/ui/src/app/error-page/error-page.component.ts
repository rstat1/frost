import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { ConfigService } from '../login/config-service';

@Component({
	selector: 'app-error-page',
	templateUrl: './error-page.html',
	styleUrls: ['./style.css']
})
export class ErrorPageComponent implements OnInit {
	public error: string = "";
	public description: string = "";

	private errors: string[] = ["Invalid Request", "Unauthorized Access", "Invalid Credentials"];
	private descriptions: string[] = [
		"Please try again later.",
		"Your account is not authorized to access this service.",
		"The username or password provided, was invalid."
	];

	constructor(private route: ActivatedRoute) {}
	ngOnInit() {
		let id = this.route.snapshot.paramMap.get('id');
		this.error = this.errors[parseInt(id)-1];
		this.description = this.descriptions[parseInt(id)-1];
		(<any>document.getElementsByClassName("background")[0]).style.backgroundImage = 'url("' + ConfigService.GetFrostURLFor("bg") + '")'
	}
}
