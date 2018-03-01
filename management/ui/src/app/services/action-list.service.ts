import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

export class PrimaryActionInfo {
	public PrimaryActionName: string;
	public PrimaryActionIcon: string;
	public PrimaryActionSubText: string;
	constructor(public name: string, public icon: string, public title: string ) {
		this.PrimaryActionName = name;
		this.PrimaryActionIcon = icon;
		this.PrimaryActionSubText = title;
	}
}

export class SubItemDetails {
	public IconName: string;
	public Description: string;
}

@Injectable()
export class ActionListService {
	public SubActionClicked: Observable<string>;
	public ActionListItems: Observable<string[]>;
	public PrimaryActionClicked: Observable<string>;
	public ClearPrimarySelection: Observable<string>;
	public HighlightPrimaryAction: Observable<string>;
	public PrimaryAction: Observable<PrimaryActionInfo>;
	public ActionListSubItems: Observable<SubItemDetails[]>;

	private clearPrimary: Subject<string>;
	private subActionClicked: Subject<string>;
	private actionListItems: Subject<string[]>;
	private primaryActionClicked: Subject<string>;
	private highlightPrimaryAction: Subject<string>;
	private primaryAction: Subject<PrimaryActionInfo>;
	private actionListSubItems: Subject<SubItemDetails[]>;

	constructor() {
		this.clearPrimary = new Subject<string>();
		this.subActionClicked = new Subject<string>();
		this.actionListItems = new Subject<string[]>();
		this.primaryActionClicked = new Subject<string>();
		this.primaryAction = new Subject<PrimaryActionInfo>();
		this.actionListSubItems = new Subject<SubItemDetails[]>();
		this.highlightPrimaryAction = new Subject<string>();

		this.PrimaryAction = this.primaryAction.asObservable();
		this.ActionListItems = this.actionListItems.asObservable();
		this.SubActionClicked = this.subActionClicked.asObservable();
		this.ClearPrimarySelection = this.clearPrimary.asObservable();
		this.ActionListSubItems = this.actionListSubItems.asObservable();
		this.PrimaryActionClicked = this.primaryActionClicked.asObservable();
		this.HighlightPrimaryAction = this.highlightPrimaryAction.asObservable();
	}
	public ClearSelectedItem() { this.clearPrimary.next("hi"); }
	public OnHighlightPrimaryAction() { this.highlightPrimaryAction.next("name"); }
	public OnSubActionClicked(subAction: string) { this.subActionClicked.next(subAction); }
	public SetActionList(newActionList: string[]) { this.actionListItems.next(newActionList); }
	public SetSubItems(newSubList: SubItemDetails[]) { this.actionListSubItems.next(newSubList); }
	public OnPrimaryActionClicked(actionName: string) { this.primaryActionClicked.next(actionName); }
	public SetPrimaryAction(newPrimaryAction: PrimaryActionInfo) { this.primaryAction.next(newPrimaryAction); }
}
