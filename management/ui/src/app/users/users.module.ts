import { Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AuthGuard } from 'app/services/auth/auth.guard';
import { NewUserComponent } from 'app/users/new/new-user';
import { UsersRootComponent } from 'app/users/root/users-root';
import { EditUserComponent } from 'app/users/edit/edit-user.component';

const projectRoutes: Routes = [
	{
		path: 'manage',
		component: UsersRootComponent,
		canActivate: [AuthGuard],
		children: [
			{ path: 'users/new', component: NewUserComponent },
			{ path: 'users/edit/:name', component: EditUserComponent },
		]
	}
];

@NgModule({
	declarations: [
		NewUserComponent,
		UsersRootComponent
	],
	imports: [
		CommonModule
	]
})
export class UsersModule { }
