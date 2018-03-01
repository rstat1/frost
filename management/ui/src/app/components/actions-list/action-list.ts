import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatTableDataSource } from '@angular/material';

import { ActionListService, SubItemDetails } from 'app/services/action-list.service';
import { Subscription } from 'rxjs/Subscription';

@Component({
	selector: 'app-action-list',
	templateUrl: './action-list.html',
	styleUrls: ['./action-list.css']
})
export class ActionsListComponent implements OnInit, OnDestroy {
	private subItemsSub: Subscription;
	private clearPrimary: Subscription;
	private actionsListSub: Subscription;
	private primaryActionIcon: string = "";
	private primaryActionName: string = "";
	private primaryActionSub: Subscription;
	private subActions: SubItemDetails[] = [];
	private highlightPrimaryAct: Subscription;
	private primaryActionDescription: string = "";
	private displayedColumns = ['Name', "Actions"];
	private dataSource = new MatTableDataSource<string>();

	public isSelected: boolean = false;

	constructor(private actionService: ActionListService) {}
	ngOnInit(): void {
		this.actionsListSub = this.actionService.ActionListItems.subscribe(items => {
			this.dataSource.data = items;
		});
		this.primaryActionSub = this.actionService.PrimaryAction.subscribe(action => {
			this.primaryActionIcon = action.PrimaryActionIcon;
			this.primaryActionName = action.PrimaryActionName;
			this.primaryActionDescription = action.PrimaryActionSubText;
		});
		this.subItemsSub = this.actionService.ActionListSubItems.subscribe(subItems => {
			this.subActions = subItems;
		});
		this.clearPrimary = this.actionService.ClearPrimarySelection.subscribe(() => {
			this.isSelected = false;
		});
		this.highlightPrimaryAct = this.actionService.HighlightPrimaryAction.subscribe(() => {
			this.isSelected = true;
		});
	}
	ngOnDestroy(): void {
		this.subItemsSub.unsubscribe();
		this.clearPrimary.unsubscribe();
		this.actionsListSub.unsubscribe();
		this.primaryActionSub.unsubscribe();
		this.highlightPrimaryAct.unsubscribe();
	}
	public primaryActionClicked(name: string) {
		this.isSelected = true;
		this.actionService.OnPrimaryActionClicked(name);
	}
	public subActionClicked(name: string) {
		this.actionService.OnSubActionClicked(name);
	}
}