import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { MatTableDataSource, MatCheckboxChange } from '@angular/material';

import { APIService } from 'app/services/api/api.service';
import { PageInfoService } from 'app/services/page-info.service';
import { PrimaryActionInfo } from 'app/services/action-list/action-list-common';
import { ActionListService } from 'app/services/action-list/action-list.service';
import { ServiceAccess, Permission, PermissionChange, PasswordChange } from 'app/services/api/api-common';
import { ConfigService } from 'app/services/config.service';

@Component({
	selector: 'app-edit-user',
	templateUrl: './edit-user.html',
	styleUrls: ['./edit-user.css']
})
export class EditUserComponent implements OnInit {
	public newPassword: string = "";
	public currentUsername: string = "";
	public dataSource = new MatTableDataSource<string>();
	public displayedColumns = ['name', 'CanAccess', 'HasRoot'];

	private permissions: ServiceAccess[] = new Array();

	constructor(private actions: ActionListService, private api: APIService, private header: PageInfoService,
		private route: ActivatedRoute) {
		this.api.GetServices(true).subscribe(resp => {
			this.dataSource.data = JSON.parse(resp.response);
		});
		this.route.params.subscribe(dumb => {
			this.api.GetUserInfo(dumb["name"]).subscribe(r => {
				let user = JSON.parse(r.response);
				this.currentUsername = user.Username;
				this.api.GetPermissionMap(this.currentUsername).subscribe(perms => {
					let m = JSON.parse(perms.response);
					this.permissions = m;
				});
			});
		});
	}
	ngOnInit() {
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New User", "add", "Add a new user"));
		this.header.SetPagePath(window.location.pathname);
	}
	public checkChanged(row: string, type: string, e: MatCheckboxChange) {
		let pc: PermissionChange = { name: type, service: row, newValue: e.checked, user: this.currentUsername };
		this.api.ChangePermissionValue(pc).subscribe(_ => { });
	}
	public getValue(service: string, permission: string): boolean {
		let p: Permission[] = this.permissions[service];
		if (p != undefined) {
			return p[permission];
		}
	}
	public changePW() {
		let pc: PasswordChange = { user: this.currentUsername, pass: this.newPassword };
		this.api.ChangePassword(pc).subscribe(_ => { });
	}
	public getServiceIconURL(name: string): string {
		return ConfigService.GetAPIURLFor("icon/" + name);
	}
}
