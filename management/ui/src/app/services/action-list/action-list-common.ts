export enum PrimaryActionState {
	Active,
	Inactive,
}
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
export interface IActionList {
	PrimaryActionClicked(name: string): void;
	SetActionListItems(items: string[]): void;
	SetUseDefaultImage(newState: boolean): void;
	SetPrimaryAction(details: PrimaryActionInfo): void;
	SubActionClicked(details: SubActionClickEvent): void;
	SetActionListSubItems(items: SubItemDetails[]): void;
	ChangePrimaryActionState(newState: PrimaryActionState): void;
}