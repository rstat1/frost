import { CommonModule } from '@angular/common';
import { NgModule, ModuleWithProviders, Optional, SkipSelf, Inject } from '@angular/core';

import { MenuItem } from "app/menu/menu-common";
import { MenuService } from 'app/services/menu.service';
import { MenuComponent } from "app/components/menu/menu.component";

export class MenuItems {
	Items: MenuItem[];
}

@NgModule({ imports: [ CommonModule ], exports: [ MenuComponent ], declarations: [ MenuComponent ], providers:[MenuItem]})
export class MenuModule {
	constructor(@Inject(MenuItems) private items: MenuItems, private menuService: MenuService) {
		if (items != null) {
			this.menuService.AddItemsToMenu(items.Items);
		}
	}
	static forRoot(items: MenuItems) : ModuleWithProviders {
		return {
			ngModule: MenuModule,
			providers: [
				{provide: MenuItems, useValue: items, deps: [MenuItem]},
			]
		}
	}
}
