import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { AuthGuard } from 'app/services/auth/auth.guard';
import { NewUserComponent } from 'app/users/new/new-user';
import { UsersRootComponent } from 'app/users/root/users-root';
import { EditUserComponent } from 'app/users/edit/edit-user.component';
import { MatTabsModule, MatIconModule, MatTableModule, MatCardModule, MatListModule,
		 MatInputModule, MatChipsModule, MatDialogModule, MatButtonModule, MatToolbarModule,
		 MatTooltipModule, MatStepperModule, MatCheckboxModule, MatSnackBarModule, MatExpansionModule,
		 MatSlideToggleModule } from '@angular/material';
import { MenuModule } from 'app/menu/menu.module';
import { UserListComponent } from './user-list/user-list.component';
import { ActionsListModule } from 'app/components/actions-list/action-list.module';

const projectRoutes: Routes = [
	{
		path: 'users',
		component: UsersRootComponent,
		canActivate: [AuthGuard],
		children: [
			{ path: 'new', component: NewUserComponent },
			{ path: 'edit/:name', component: EditUserComponent },
			{ path: '', component: UserListComponent },
		]
	}
];

@NgModule({
	declarations: [
		NewUserComponent,
		UserListComponent,
		EditUserComponent,
		UsersRootComponent,
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
		RouterModule.forChild(projectRoutes)
	]
})
export class UsersModule { }
