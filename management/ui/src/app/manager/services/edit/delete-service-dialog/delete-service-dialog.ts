import { Component, Inject } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';

import { MenuService } from 'app/services/menu.service';
import { APIService } from 'app/services/api/api.service';

@Component({
	selector: 'app-delete-service-dialog',
	templateUrl: './delete-service-dialog.html',
	styleUrls: ['./delete-service-dialog.css']
})
export class DeleteServiceDialogComponent {
	public serviceName: string;
	public typedServiceName: string;
	public enteredNameProperly: boolean = true;
	public errorMessage: string = "Please enter the service name.";

	constructor(@Inject(MAT_DIALOG_DATA) public data: any, public router: Router, private route: ActivatedRoute,
				public dialogRef: MatDialogRef<DeleteServiceDialogComponent>, private api: APIService,
				private menu: MenuService) {
		this.serviceName = <string>data.project;
	}
	public save() {
		if (this.typedServiceName == this.serviceName) {
			this.enteredNameProperly = true;
			if (this.serviceName != "watchdog") {
				this.api.DeleteService(this.serviceName).subscribe(resp => {
					if (resp.status == "success") {
						this.dialogRef.close(true);
						this.router.navigate(['services']);
						// this.menu.SetMenuContext("services", "");
					}
				}, error => {
					this.enteredNameProperly = false;
					this.errorMessage = error;
				});
			}
		} else {
			this.enteredNameProperly = false;
		}
	}
}
