import { Subject, Observable } from 'rxjs';
import { Injectable } from '@angular/core';
import { PrimaryActionInfo, SubItemDetails, SubActionClickEvent, IActionList, PrimaryActionState } from 'app/services/action-list/action-list-common';

class NullActionList implements IActionList {
	PrimaryActionClicked(name: string): void {}
	SetActionListItems(items: string[]): void {}
	SetUseDefaultImage(newState: boolean): void {}
	SetPrimaryAction(details: PrimaryActionInfo): void {}
	SubActionClicked(details: SubActionClickEvent): void {}
	SetActionListSubItems(items: SubItemDetails[]): void {}
	ChangePrimaryActionState(newState: PrimaryActionState): void {}
}

@Injectable()
export class ActionListService {
	public PrimaryActionClicked: Observable<string>;
	public SubActionClicked: Observable<SubActionClickEvent>;

	private primaryActionClicked: Subject<string>;
	private subActionClicked: Subject<SubActionClickEvent>;

	private list: IActionList;

	constructor() {
		this.list = new NullActionList();
		this.primaryActionClicked = new Subject<string>();
		this.subActionClicked = new Subject<SubActionClickEvent>();

		this.SubActionClicked = this.subActionClicked.asObservable();
		this.PrimaryActionClicked = this.primaryActionClicked.asObservable();
	}
	public SetActionListInstance(actionList: IActionList) {
		this.list = actionList;
	}
	public ClearSelectedItem() {
		this.list.ChangePrimaryActionState(PrimaryActionState.Inactive); }
	public OnHighlightPrimaryAction() { this.list.ChangePrimaryActionState(PrimaryActionState.Active); }
	public SetActionList(newActionList: string[]) { this.list.SetActionListItems(newActionList); }
	public SetSubItems(newSubList: SubItemDetails[]) { this.list.SetActionListSubItems(newSubList); }
	public OnPrimaryActionClicked(actionName: string) { this.primaryActionClicked.next(actionName); }
	public SetPrimaryAction(newPrimaryAction: PrimaryActionInfo) { this.list.SetPrimaryAction(newPrimaryAction); }
	public OnSubActionClicked(subActionName: string, ctx: string) {
		this.subActionClicked.next({ContextInfo: ctx, SubActionName: subActionName});
	}
	public SetImageType(useDefault: boolean) { this.list.SetUseDefaultImage(useDefault); }
}
