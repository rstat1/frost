import { MatTableDataSource } from '@angular/material';
import { Component, OnInit, OnDestroy } from '@angular/core';

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
	private primaryActionSub: Subscription;
	private useDefaultImageSub: Subscription;
	private highlightPrimaryAct: Subscription;

	public showAction: boolean;
	public useDefaultImage: boolean;
	public isSelected: boolean = false;
	public primaryActionIcon: string = "";
	public primaryActionName: string = "";
	public subActions: SubItemDetails[] = [];
	public primaryActionDescription: string = "";
	public displayedColumns = ['Name', "Actions"];
	public dataSource = new MatTableDataSource<string>();

	constructor(private actionService: ActionListService) {}
	ngOnInit(): void {
		this.useDefaultImageSub = this.actionService.UseDefaultImage.subscribe(img => {
			console.log(img);
			this.useDefaultImage = img;
		});
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
		this.useDefaultImageSub.unsubscribe();
		this.highlightPrimaryAct.unsubscribe();
	}
	public primaryActionClicked(name: string) {
		this.isSelected = true;
		this.actionService.OnPrimaryActionClicked(name);
	}
	public subActionClicked(name: string, ctx: string) {
		this.actionService.OnSubActionClicked(name, ctx);
	}
}