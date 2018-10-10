import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule, MatButtonModule, MatCheckboxModule, MatTableModule,
		MatToolbarModule, MatIconModule, MatSnackBarModule, MatTooltipModule,
		MatSlideToggleModule, MatExpansionModule, MatListModule, MatDialogModule } from '@angular/material';
import { FormsModule, FormBuilder, ReactiveFormsModule } from '@angular/forms';

import { MenuModule } from 'app/menu/menu.module';
import { AuthGuard } from 'app/services/auth/auth.guard';
import { EditServiceComponent } from './services/edit/edit';
import { PageInfoService } from 'app/services/page-info.service';
import { NewUserComponent } from 'app/manager/users/new/new-user';
import { ActionListService } from 'app/services/action-list.service';
import { UsersRootComponent } from 'app/manager/users/root/users-root';
import { NewServiceComponent } from 'app/manager/services/new/new-service';
import { ManagerRootComponent } from 'app/manager/manager-root/manager-root';
import { ActionsListComponent } from 'app/components/actions-list/action-list';
import { EditUserComponent } from 'app/manager/users/edit/edit-user.component';
import { ServicesRootComponent } from 'app/manager/services/root/services-root';
import { LogViewerComponent } from 'app/manager/log-viewer/log-viewer.component';
import { NewAliasDialogComponent } from 'app/manager/services/edit/new-alias-dialog/new-alias';

const projectRoutes: Routes = [
	{
		path: 'manage',
		component: ManagerRootComponent,
		canActivate: [AuthGuard],
		children: [
			{ path: 'logs', component: LogViewerComponent },
			{ path: 'users', component: UsersRootComponent },
			{ path: 'users/new', component: NewUserComponent },
			{ path: 'users/edit/:name', component: EditUserComponent },
			{ path: 'services', component: ServicesRootComponent },
			{ path: 'services/new', component: NewServiceComponent },
			{ path: 'services/edit/:name', component: EditServiceComponent }
		]
	}
];
@NgModule({
	imports: [
		FormsModule,
		CommonModule,
		MatListModule,
		MatIconModule,
		MatTableModule,
		MatInputModule,
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
		MenuModule.forRoot(null),
		RouterModule.forChild(projectRoutes)
	],
	declarations: [
		NewUserComponent,
		UsersRootComponent,
		LogViewerComponent,
		NewServiceComponent,
		ManagerRootComponent,
		ActionsListComponent,
		EditServiceComponent,
		ServicesRootComponent,
		NewAliasDialogComponent,
		EditUserComponent,
	],
	providers: [ FormBuilder, ActionListService, PageInfoService ],
	entryComponents: [ NewAliasDialogComponent ],
})
export class ManagerModule { }
