import { Subscription } from 'rxjs';
import { MatTableDataSource } from '@angular/material';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { environment } from 'environments/environment';
import { ActionListService } from 'app/services/action-list/action-list.service';
import { SubItemDetails, IActionList, SubActionClickEvent, PrimaryActionState, PrimaryActionInfo } from 'app/services/action-list/action-list-common';
import { ConfigService } from 'app/services/config.service';

@Component({
	selector: 'app-action-list',
	templateUrl: './action-list.html',
	styleUrls: ['./action-list.css']
})
export class ActionsListComponent implements OnInit, OnDestroy, IActionList {
	public showAction: boolean;
	public useDefaultImage: boolean;
	public isSelected: boolean = false;
	public primaryActionIcon: string = "";
	public primaryActionName: string = "";
	public subActions: SubItemDetails[] = [];
	public primaryActionDescription: string = "";
	public displayedColumns = ['Name', "Actions"];
	public dataSource = new MatTableDataSource<string>();

	constructor(private actionService: ActionListService) { this.actionService.SetActionListInstance(this); }
	ngOnInit(): void {
		console.log("action list onInit");
	}
	ngOnDestroy(): void {
	}
	PrimaryActionClicked(): void {
		this.isSelected = true;
		this.actionService.OnPrimaryActionClicked(name);
	}
	SetActionListItems(items: string[]): void {
		this.dataSource.data = items;
	}
	SubActionClicked(details: SubActionClickEvent): void {
		this.actionService.OnSubActionClicked(details.SubActionName, details.ContextInfo);
	}
	SetActionListSubItems(items: SubItemDetails[]): void {
		this.subActions = items;
	}
	ChangePrimaryActionState(newState: PrimaryActionState): void {
		this.isSelected = newState == PrimaryActionState.Active;
	}
	SetUseDefaultImage(newState: boolean): void {
		this.useDefaultImage = newState;
	}
	SetPrimaryAction(details: PrimaryActionInfo): void {
		this.primaryActionIcon = details.PrimaryActionIcon;
		this.primaryActionName = details.PrimaryActionName;
		this.primaryActionDescription = details.PrimaryActionSubText;
	}
	public primaryActionClicked(name: string) {
		this.isSelected = true;
		this.actionService.OnPrimaryActionClicked(name);
	}
	public subActionClicked(name: string, ctx: string) {
		this.actionService.OnSubActionClicked(name, ctx);
	}
	public getServiceIconURL(name: string): string {
		return ConfigService.GetAPIURLFor("icon/" + name);
	}
}