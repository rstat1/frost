import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule, MatButtonModule, MatCheckboxModule, MatTableModule,
	 	 MatToolbarModule, MatIconModule, MatSnackBarModule, MatTooltipModule } from '@angular/material';
import { FormsModule, FormGroup, FormBuilder, ReactiveFormsModule } from '@angular/forms';

import { MenuModule } from 'app/menu/menu.module';
import { UsersRootComponent } from './users/root/users-root';
import { PageInfoService } from 'app/services/page-info.service';
import { NewServiceComponent } from './services/new/new-service';
import { NewUserComponent } from './users/new/new-user.component';
import { MenuComponent } from 'app/components/menu/menu.component';
import { ManagerRootComponent } from './manager-root/manager-root';
import { ActionListService } from '../services/action-list.service';
import { AuthGuard, RootGuard } from 'app/services/auth/auth.guard';
import { ServicesRootComponent } from './services/root/services-root';
import { LogViewerComponent } from './log-viewer/log-viewer.component';
import { ActionsListComponent } from 'app/components/actions-list/action-list';

const projectRoutes: Routes = [
	{
		path: 'manage',
		component: ManagerRootComponent,
		canActivate: [AuthGuard],
		children: [
			{ path: 'logs', component: LogViewerComponent},
			{ path: 'users', component: UsersRootComponent},
			{ path: 'users/new', component: NewUserComponent},
			{ path: 'services', component: ServicesRootComponent},
			{ path: 'services/new', component: NewServiceComponent}
		]
	}
]
@NgModule({
	imports: [
		FormsModule,
		CommonModule,
		MatIconModule,
		MatTableModule,
		MatInputModule,
		MatButtonModule,
		MatToolbarModule,
		MatTooltipModule,
		MatStepperModule,
		MatCheckboxModule,
		MatSnackBarModule,
		ReactiveFormsModule,
		MenuModule.forRoot(null),
		// MalihuScrollbarModule.forRoot(),
		RouterModule.forChild(projectRoutes)
  	],
  	declarations: [
		NewUserComponent,
		UsersRootComponent,
		LogViewerComponent,
		ManagerRootComponent,
		NewServiceComponent,
		ActionsListComponent,
		ServicesRootComponent,
	],
	providers: [ FormBuilder, ActionListService, PageInfoService ],
	// entryComponents: [ ProjectListItem ],
})
export class ManagerModule { }
