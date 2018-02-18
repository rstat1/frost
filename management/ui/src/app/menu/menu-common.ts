export class MenuItem {
	Icon: string;
	Category?: string;
	Context?: string;
	ItemTitle: string;
	ActionName: string;
	ItemSubtext: string;
	RequiresRoot?: boolean;
	SecondaryContext?: string;
}
export class Context {
	Extra: string;
	ContextName: string;
}