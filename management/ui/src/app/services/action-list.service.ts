import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

export class SubItemDetails {
	public IconName: string;
	public Description: string;
}
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
export class SubActionClickEvent {
	public ContextInfo: string;
	public SubActionName: string;
}
@Injectable()
export class ActionListService {
	public UseDefaultImage: Observable<boolean>;
	public ActionListItems: Observable<string[]>;
	public PrimaryActionClicked: Observable<string>;
	public ClearPrimarySelection: Observable<string>;
	public HighlightPrimaryAction: Observable<string>;
	public PrimaryAction: Observable<PrimaryActionInfo>;
	public ActionListSubItems: Observable<SubItemDetails[]>;
	public SubActionClicked: Observable<SubActionClickEvent>;

	private clearPrimary: Subject<string>;
	private useDefaultImg: Subject<boolean>;
	private actionListItems: Subject<string[]>;
	private primaryActionClicked: Subject<string>;
	private highlightPrimaryAction: Subject<string>;
	private primaryAction: Subject<PrimaryActionInfo>;
	private actionListSubItems: Subject<SubItemDetails[]>;
	private subActionClicked: Subject<SubActionClickEvent>;

	constructor() {
		this.clearPrimary = new Subject<string>();
		this.useDefaultImg = new Subject<boolean>();
		this.actionListItems = new Subject<string[]>();
		this.primaryActionClicked = new Subject<string>();
		this.highlightPrimaryAction = new Subject<string>();
		this.primaryAction = new Subject<PrimaryActionInfo>();
		this.actionListSubItems = new Subject<SubItemDetails[]>();
		this.subActionClicked = new Subject<SubActionClickEvent>();

		this.PrimaryAction = this.primaryAction.asObservable();
		this.UseDefaultImage = this.useDefaultImg.asObservable();
		this.ActionListItems = this.actionListItems.asObservable();
		this.SubActionClicked = this.subActionClicked.asObservable();
		this.ClearPrimarySelection = this.clearPrimary.asObservable();
		this.ActionListSubItems = this.actionListSubItems.asObservable();
		this.PrimaryActionClicked = this.primaryActionClicked.asObservable();
		this.HighlightPrimaryAction = this.highlightPrimaryAction.asObservable();


	}
	public ClearSelectedItem() { this.clearPrimary.next("hi"); }
	public OnHighlightPrimaryAction() { this.highlightPrimaryAction.next("name"); }
	public SetActionList(newActionList: string[]) { this.actionListItems.next(newActionList); }
	public SetSubItems(newSubList: SubItemDetails[]) { this.actionListSubItems.next(newSubList); }
	public OnPrimaryActionClicked(actionName: string) { this.primaryActionClicked.next(actionName); }
	public SetPrimaryAction(newPrimaryAction: PrimaryActionInfo) { this.primaryAction.next(newPrimaryAction); }
	public OnSubActionClicked(subActionName: string, ctx: string) {
		this.subActionClicked.next({ContextInfo: ctx, SubActionName: subActionName});
	}
	public SetImageType(useDefault: boolean) { this.useDefaultImg.next(useDefault); }
}
