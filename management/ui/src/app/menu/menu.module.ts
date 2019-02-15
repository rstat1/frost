import { CommonModule } from '@angular/common';
import { MatTooltipModule } from '@angular/material';
import { NgModule, ModuleWithProviders, Inject } from '@angular/core';

import { MenuService } from 'app/services/menu.service';
import { MalihuScrollbarModule } from 'ngx-malihu-scrollbar';
import { MenuComponent } from 'app/components/menu/menu.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

export class MenuItem {
	Icon: string;
	Category: string;
	Context?: string;
	MenuType?: string;
	ItemTitle: string;
	ActionName: string;
	ItemSubtext: string;
	RequiresRoot?: boolean;
}
export class MenuItems {
	Items: MenuItem[];
}

@NgModule({
	imports: [ CommonModule, MalihuScrollbarModule.forRoot(), MatTooltipModule, BrowserAnimationsModule ],
	exports: [ MenuComponent ],
	declarations: [ MenuComponent ],
	providers:[ MenuItem ]
})
export class MenuModule {
	constructor(@Inject(MenuItems) private items: MenuItems[], private menuService: MenuService) {
		if (items.length >= 1) {
			items.forEach(list => {
				this.menuService.AddItemsToMenu(list.Items);
			});
		}
		// if (items != null) {
		// }
	}
	static forRoot(items: MenuItems): ModuleWithProviders {
		return {
			ngModule: MenuModule,
			providers: [
				{provide: MenuItems, multi: true, useValue: items, deps: [MenuItem]},
			]
		};
	}
	static forChild(items: MenuItems): ModuleWithProviders<MenuModule> {
		return {ngModule: MenuModule, providers: [
			{ provide: MenuItems, multi: true, useValue: items, deps: [ MenuItem ] },
		]};
	}
}