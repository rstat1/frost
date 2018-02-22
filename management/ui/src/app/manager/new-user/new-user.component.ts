import { Component, OnInit } from '@angular/core';
import {MatTableDataSource} from '@angular/material';
import {SelectionModel} from '@angular/cdk/collections';

import { AuthRequest } from 'app/services/api/api.service';
import { ServicePermission } from 'app/services/api/api-common';

const DATA: ServicePermission[] = [
	{name: "gemini", hasRoot: false, hasAccess: false},
	{name: "player3", hasRoot: false, hasAccess: false},
	{name: "watchdog", hasRoot: false, hasAccess: false},
]

@Component({
	selector: 'app-new-user',
	templateUrl: './new-user.html',
	styleUrls: ['./new-user.css']
})
export class NewUserComponent implements OnInit {
	private m: AuthRequest = new AuthRequest("","");
	private dataSource = new MatTableDataSource<ServicePermission>(DATA);
	private displayedColumns = ['name', 'CanAccess', 'HasRoot'];

	constructor() { }
	ngOnInit() { }
	public save() {}
}
