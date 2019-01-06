import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatTableDataSource, MatSnackBar } from '@angular/material';

import { APIService } from 'app/services/api/api.service';
import { PageInfoService } from 'app/services/page-info.service';
import { AuthRequest, NewUser, ServiceAccess } from 'app/services/api/api-common';
import { ActionListService, PrimaryActionInfo } from 'app/services/action-list.service';
import { environment } from 'environments/environment';

@Component({
	selector: 'app-new-user',
	templateUrl: './new-user.html',
	styleUrls: ['./new-user.css']
})
export class NewUserComponent implements OnInit, OnDestroy {
	private status: string = "";
	public m: AuthRequest = new AuthRequest("", "");
	public dataSource = new MatTableDataSource<string>();
	public displayedColumns = ['name', 'CanAccess', 'HasRoot'];
	private selection = new SelectionModel<string>(true, []);
	private permissions: ServiceAccess[] = new Array();

	constructor(private api: APIService, private snackBar: MatSnackBar,
		private actions: ActionListService, private header: PageInfoService) {
		this.api.GetServices(true).subscribe(resp => {
			this.dataSource.data = JSON.parse(resp.response);
		});
	}
	ngOnInit() {
		console.warn("NewUserComponent init...");
		this.actions.OnHighlightPrimaryAction();
		this.actions.SetPrimaryAction(new PrimaryActionInfo("New User", "add", "Add a new user"));
		this.header.SetPagePath(window.location.pathname);
	}
	ngOnDestroy() {
		console.warn("NewUserComponent goes away");
		this.actions.ClearSelectedItem();
	}
	public save() {
		const user: NewUser = new NewUser(this.m.Username, this.m.Password, this.permissions);
		this.api.SaveUser(user).subscribe(
			resp => {
				this.snackBar.open("Successfully added new user", "", {
					duration: 3000,
					panelClass: "proper-colors",
					horizontalPosition: 'right',
					verticalPosition: 'top',
				});
			},
			err => {
				this.snackBar.open(err.error.response, "", {
					duration: 3000,
					panelClass: "proper-colors",
					horizontalPosition: 'right',
					verticalPosition: 'top',
				});
			}
		);
	}
	public getServiceIconURL(name: string): string {
		return environment.APIBaseURL + "/frost/icon/"+name;
	}
	private checkChanged(row: any, type: string, checkEvent: any) {
		const service = this.permissions.find(s => s.service == row);
		if (service == undefined) {
			this.permissions.push({service: row, permissions: [{name: type, value: true}]});
		} else {
			const permission = service.permissions.find(p => p.name == type);
			if (permission == undefined) { service.permissions.push({name: type, value: true}); }
			else { permission.value = !permission.value; }
		}
	}
}