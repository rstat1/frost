import { Subject, ReplaySubject, Observable } from 'rxjs';
import { Injectable } from '@angular/core';

import { MenuItem } from 'app/menu/menu.module';
import { AuthService } from 'app/services/auth/auth.service';

@Injectable({ providedIn: 'root' })
export class MenuService {
	public MenuItemClicked: Observable<string>;
	public CategoryChanged: Observable<string>;

	private menuItemClicked: Subject<string>;
	private categoryChanged: ReplaySubject<string>;
	public categories: string[] = new Array();
	private menuContext: Map<string, string> = new Map();
	private menuItemToCategory: Map<string, Array<MenuItem>> = new Map();

	constructor(private auth: AuthService) {
		this.menuItemClicked = new Subject<string>();
		this.categoryChanged = new ReplaySubject<string>();

		this.MenuItemClicked = this.menuItemClicked.asObservable();
		this.CategoryChanged = this.categoryChanged.asObservable();
	}
	public AddItemsToMenu(items: MenuItem[]) {
		console.log("add items");
		if (items != null) {
			if (this.auth.IsLoggedIn == false) {
				this.auth.TokenValidation.subscribe(result => {
					if (result) { this.loadMenuItems(items); }
				});
			} else {
				this.loadMenuItems(items);
			}
		}
	}
	public GetMenuContextData(name: string): string {
		return this.menuContext.get(name);
	}
	public SetMenuContext(currentPage: string, extra: string) {
		this.menuContext.set("currentPage", currentPage);
		this.menuContext.set("extra", extra);
	}
	public SetMenuCategory(newCategory: string) {
		this.categoryChanged.next(newCategory);
	}
	public HandleMouseEvent(clickedItemTitle: string) {
		this.menuItemClicked.next(clickedItemTitle);
	}
	public GetCategoryItems(category: string, menuType?: string): MenuItem[] {
		let currentPage: string = this.GetMenuContextData("currentPage");
		let items: MenuItem[] = this.menuItemToCategory.get(category);
		if (items != null) {
			let contextualItems: MenuItem[] = items.filter((item) => {
				if (this.checkMenuType(menuType, item.MenuType) == true) {
					if (item.Context != null) {
						if (item.Context.startsWith("!") && item.Context.substring(1, item.Context.length) != currentPage) {
							return item;
						} else if (item.Context == currentPage) {
							return item;
						}
					} else {
						return item;
					}
				}
			});
			return contextualItems;
		} else {
			return null;
		}
	}
	public GetCategoryList(): string[] {
		return this.categories;
	}
	private checkMenuType(menuType: string, itemMenuType: string): boolean {
		return ((menuType != null && itemMenuType != null) && menuType == itemMenuType);
	}
	private isRootItemAllowed(item: MenuItem): boolean {
		return item.RequiresRoot && this.auth.UserIsRoot;
	}
	private loadMenuItems(items: MenuItem[]) {
		if (this.categories.length == 0) {
			items.forEach(item => {
				if (this.menuItemToCategory.has(item.Category) == false) {
					if (this.isRootItemAllowed(item) == true) {
						this.menuItemToCategory.set(item.Category, [ item ]);
						this.categories.push(item.Category);
					}
					else if (item.RequiresRoot == false) {
						this.menuItemToCategory.set(item.Category, [ item ]);
						this.categories.push(item.Category);
					}
				} else {
					let cat: MenuItem[] = this.menuItemToCategory.get(item.Category);
					if (this.isRootItemAllowed(item)) { cat.push(item); }
					else if (item.RequiresRoot == false || item.RequiresRoot == undefined) { cat.push(item); }
					this.menuItemToCategory.set(item.Category, cat);
				}
			});
		}
	}
}