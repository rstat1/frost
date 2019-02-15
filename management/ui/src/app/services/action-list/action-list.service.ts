import { Subject, Observable } from 'rxjs';
import { Injectable } from '@angular/core';
import { PrimaryActionInfo, SubItemDetails, SubActionClickEvent, IActionList, PrimaryActionState } from 'app/services/action-list/action-list-common';


@Injectable()
export class ActionListService {
	public PrimaryActionClicked: Observable<string>;
	public SubActionClicked: Observable<SubActionClickEvent>;

	private primaryActionClicked: Subject<string>;
	private subActionClicked: Subject<SubActionClickEvent>;

	private list: IActionList;

	constructor() {
		this.primaryActionClicked = new Subject<string>();
		this.subActionClicked = new Subject<SubActionClickEvent>();

		this.SubActionClicked = this.subActionClicked.asObservable();
		this.PrimaryActionClicked = this.primaryActionClicked.asObservable();
	}
	public SetActionListInstance(actionList: IActionList) {
		this.list = actionList;
	}
	public ClearSelectedItem() { this.list.ChangePrimaryActionState(PrimaryActionState.Inactive); }
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
