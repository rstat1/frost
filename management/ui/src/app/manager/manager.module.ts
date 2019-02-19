import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormsModule } from '@angular/forms';
import { MatListModule, MatIconModule, MatTableModule,
	MatInputModule, MatChipsModule, MatDialogModule,
	MatButtonModule, MatToolbarModule, MatTooltipModule,
	MatStepperModule, MatCheckboxModule, MatSnackBarModule,
	MatExpansionModule, MatSlideToggleModule, MatTabsModule, MatCardModule } from '@angular/material';

import { MenuModule } from 'app/menu/menu.module';
import { ActionsListModule } from 'app/components/actions-list/action-list.module';

import { AuthGuard } from 'app/services/auth/auth.guard';
import { PageInfoService } from 'app/services/page-info.service';
import { ChartComponent } from 'app/components/chart/chart.component';
import { EditServiceComponent } from 'app/manager/services/edit/edit';
import { NewServiceComponent } from 'app/manager/services/new/new-service';
import { ManagerRootComponent } from 'app/manager/root/manager-root.component';
import { ActionListService } from 'app/services/action-list/action-list.service';
import { ServiceListComponent } from 'app/manager/services/list/service-list.component';

const managerRoutes: Routes = [
	{
		path: 'services',
		component: ManagerRootComponent,
		canActivate: [AuthGuard],
		children: [
			{ path: 'new', component: NewServiceComponent, pathMatch: "full" },
			{ path: '', component: ServiceListComponent, pathMatch: "full"},
			{ path: ':name', component: EditServiceComponent, pathMatch: "full" },
		]
	},
];

@NgModule({
	declarations: [
		ChartComponent,
		NewServiceComponent,
		ManagerRootComponent,
		EditServiceComponent,
		ServiceListComponent,
	],
	imports: [
		MenuModule,
		FormsModule,
		CommonModule,
		MatCardModule,
		MatListModule,
		MatTabsModule,
		MatIconModule,
		MatTableModule,
		MatInputModule,
		MatChipsModule,
		MatDialogModule,
		MatButtonModule,
		MatToolbarModule,
		MatTooltipModule,
		MatStepperModule,
		ActionsListModule,
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
