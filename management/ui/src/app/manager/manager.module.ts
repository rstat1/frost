import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule, MatButtonModule, MatCheckboxModule } from '@angular/material';
import { FormsModule, FormGroup, FormBuilder, ReactiveFormsModule } from '@angular/forms';

import { MenuModule } from 'app/menu/menu.module';
import { MenuComponent } from 'app/components/menu/menu.component';
import { AuthGuard, RootGuard } from 'app/services/auth/auth.guard';
import { ManagerRootComponent } from './manager-root/manager-root.component';
import { LogViewerComponent } from './log-viewer/log-viewer.component';

const projectRoutes: Routes = [
	{
		path: 'manage',
		component: ManagerRootComponent,
		// canActivate: [AuthGuard],
		// canActivateChild: [AuthGuard],
		children: [
			{ path: 'logs', component: LogViewerComponent}
		]
	}
]
@NgModule({
	imports: [
		FormsModule,
		CommonModule,
		// TreeViewModule,
		MatInputModule,
		MatButtonModule,
		MatStepperModule,
		MatCheckboxModule,
		ReactiveFormsModule,
		MenuModule.forRoot(null),
		// MalihuScrollbarModule.forRoot(),
		RouterModule.forChild(projectRoutes)
  	],
  	declarations: [
		LogViewerComponent,
		ManagerRootComponent,
	],
	providers: [ FormBuilder ],
	// entryComponents: [ ProjectListItem ],
})
export class ManagerModule { }
