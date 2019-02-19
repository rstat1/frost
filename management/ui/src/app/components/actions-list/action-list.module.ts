import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule, MatTooltipModule, MatIconModule, MatButtonModule } from '@angular/material';

import { ActionsListComponent } from 'app/components/actions-list/action-list';

@NgModule({
	imports: [
		CommonModule,
		MatIconModule,
		MatTableModule,
		MatButtonModule,
		MatTooltipModule,
	],
	exports: [ ActionsListComponent ],
	declarations: [ ActionsListComponent ],
})
export class ActionsListModule {}