import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { MenuItem, Context } from 'app/menu/menu-common';

@Injectable()
export class MenuService {
	public ContextChanged: Observable<Context>;
	public MenuItemClicked: Observable<string>;
	public SecondaryContextChanged: Observable<string>;

	private contextChange: Subject<Context>;
	private secondaryContextChanged: Subject<string>;

	private secondary: string = "";
	private menuItemClicked: Subject<string>;
	private menuItems: MenuItem[] = new Array();
	private menuContext: Map<string, string> = new Map();

	constructor() {
		this.contextChange = new Subject<Context>();
		this.menuItemClicked = new Subject<string>();
		this.secondaryContextChanged = new Subject<string>();

		this.ContextChanged = this.contextChange.asObservable();
		this.MenuItemClicked = this.menuItemClicked.asObservable();
		this.SecondaryContextChanged = this.secondaryContextChanged.asObservable();
	}
	public AddItemsToMenu(items: MenuItem[]) {
		this.menuItems = this.menuItems.concat(items);
	}
	public GetMenuItems(): MenuItem[] { return this.menuItems; }
	public GetMenuContextData(name: string): string {
		return this.menuContext.get(name);
	}
	public GetSecondaryCtx(): string {
		return this.secondary;
	}
	public SetMenuContext(currentPage: string, extra: string) {
		this.menuContext.set("currentPage", currentPage);
		this.menuContext.set("extra", extra);
		this.contextChange.next({ContextName: currentPage, Extra: extra});
	}
	public SetSecondaryContext(ctx: string) {
		this.secondary = ctx;
		this.secondaryContextChanged.next(this.secondary);
	}
	public HandleMouseEvent(clickedItemTitle: string){
		this.menuItemClicked.next(clickedItemTitle);
	}
}
