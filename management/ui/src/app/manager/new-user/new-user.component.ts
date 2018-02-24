import { MatTableDataSource } from '@angular/material';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnInit, Renderer2 } from '@angular/core';

import { APIService } from 'app/services/api/api.service';
import { AuthRequest, NewUser, ServiceAccess } from 'app/services/api/api-common';

@Component({
	selector: 'app-new-user',
	templateUrl: './new-user.html',
	styleUrls: ['./new-user.css']
})
export class NewUserComponent implements OnInit {
	private m: AuthRequest = new AuthRequest("","");
	private dataSource = new MatTableDataSource<string>();
	private displayedColumns = ['name', 'CanAccess', 'HasRoot'];
	private selection = new SelectionModel<string>(true, []);
	private permissions: ServiceAccess[] = new Array();

	ngOnInit() { }
	constructor(private api: APIService) {
		this.api.GetServices(true).subscribe(resp => {
			this.dataSource.data = JSON.parse(resp.response)
		});
	}
	public save() {
		let user: NewUser = new NewUser(this.m.Username, this.m.Password, this.permissions)
		console.log(JSON.stringify(user));
		this.api.SaveUser(user).subscribe();
	}
	private checkChanged(row: any, type: string, checkEvent: any) {
		let service = this.permissions.find(s => s.service == row);
		if (service == undefined) {
			this.permissions.push({service: row, permissions: [{name: type, value: true}]})
		}
		else {
			let permission = service.permissions.find(p => p.name == type)
			if (permission == undefined) { service.permissions.push({name: type, value: true}); }
			else { permission.value = !permission.value; }
		}
		console.log(this.permissions)
	}
}