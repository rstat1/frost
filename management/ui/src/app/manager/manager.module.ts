import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormsModule } from '@angular/forms';
import { MatListModule, MatIconModule, MatTableModule,
	MatInputModule, MatChipsModule, MatDialogModule,
	MatButtonModule, MatToolbarModule, MatTooltipModule,
	MatStepperModule, MatCheckboxModule, MatSnackBarModule,
	MatExpansionModule, MatSlideToggleModule } from '@angular/material';

import { MenuModule } from 'app/menu/menu.module';
import { AuthGuard } from 'app/services/auth/auth.guard';
import { PageInfoService } from 'app/services/page-info.service';
import { EditServiceComponent } from 'app/manager/services/edit/edit';
import { NewServiceComponent } from 'app/manager/services/new/new-service';
import { ManagerRootComponent } from 'app/manager/root/manager-root.component';
import { ActionsListComponent } from 'app/components/actions-list/action-list';
import { ActionListService } from 'app/services/action-list/action-list.service';

const managerRoutes: Routes = [
	{
		path: 'services',
		component: ManagerRootComponent,
		canActivate: [AuthGuard],
		children: [
			{ path: 'new', component: NewServiceComponent, pathMatch: "full" },
			{ path: ':name', component: EditServiceComponent, pathMatch: "full" },
		]
	},
];

@NgModule({
	declarations: [
		NewServiceComponent,
		ManagerRootComponent,
		ActionsListComponent,
		EditServiceComponent,
	],
	imports: [
		MenuModule,
		FormsModule,
		CommonModule,
		MatListModule,
		MatIconModule,
		MatTableModule,
		MatInputModule,
		MatChipsModule,
		MatDialogModule,
		MatButtonModule,
		MatToolbarModule,
		MatTooltipModule,
		MatStepperModule,
		MatCheckboxModule,
		MatSnackBarModule,
		MatExpansionModule,
		ReactiveFormsModule,
		MatSlideToggleModule,
		RouterModule.forChild(managerRoutes)
	],
	providers: [ ActionListService, FormBuilder, PageInfoService ]
})
export class ManagerModule { }
