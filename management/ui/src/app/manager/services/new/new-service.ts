import { NgForm } from '@angular/forms';
import { MatSnackBar } from '@angular/material';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, AbstractControl, FormControl, ValidationErrors } from '@angular/forms';

import { Service } from 'app/services/api/api-common';
import { PageInfoService } from 'app/services/page-info.service';
import { ActionListService } from 'app/services/action-list.service';
import { APIService } from 'app/services/api/api.service';

@Component({
	selector: 'app-new-service',
	templateUrl: './new-service.html',
	styleUrls: ['./new-service.css']
})
export class NewServiceComponent implements OnInit {
	public uiFiles: File;
	public serviceBin: File;
	public uiFilesName: string = "";
	public s: Service = new Service();
	public files: FormGroup;
	public serviceDetails: FormGroup;
	public managementDetails: FormGroup;

	constructor(private actions: ActionListService, private header: PageInfoService,
		private snackBar: MatSnackBar, private api: APIService, private formBuilder: FormBuilder) {}

	ngOnInit() {
		this.actions.OnHighlightPrimaryAction();
		this.header.SetPagePath(window.location.pathname);
		this.serviceDetails = this.formBuilder.group({
			'ServiceName': new FormControl('', [
				this.ValidateServiceName,
			]),
			'apiPrefix': new FormControl('', []),
			'address': new FormControl('', []),
			'authCallback': new FormControl('', []),
		});
		this.managementDetails = this.formBuilder.group({
			'IsManaged': new FormControl(false, []),
			'UpdatesManaged': new FormControl(false, [])
		});
	}
	public setFile(name: string, event: any) {
		if (name == "ui") {
			this.uiFiles = event.target.files[0];
			if (this.uiFiles.type != "application/zip") {
				this.uiFiles = null;
				this.uiFilesName = "";
				this.snackBar.open("That wasn't a zip file. >_>", "", {
					duration: 3000, panelClass: "proper-colors", horizontalPosition: 'center',
					verticalPosition: 'top',
				});
			} else {
				this.uiFilesName = this.uiFiles.name;
			}
		} else {
			this.serviceBin = event.target.files[0];
			this.s.filename = this.serviceBin.name;
		}
	}
	public save(form: NgForm) {
		let uiReader: FileReader;
		let serviceReader: FileReader;
		let body: FormData = new FormData();
		let serviceDetails: string = JSON.stringify(this.s);

		body.append("details", serviceDetails);

		if (this.uiFiles != null) {
			// console.log("have ui files...")
			body.append("uiblob", this.uiFiles, this.uiFiles.name);
		}
		if (this.serviceBin != null) {
			// console.log("have serviceBin files...")
			body.append("service", this.serviceBin, this.serviceBin.name);
		}
		this.api.NewService(body).subscribe(
			resp => {
				this.snackBar.open("Successfully added new service", "", {
					duration: 3000, panelClass: "proper-colors", horizontalPosition: 'center',
					verticalPosition: 'top',
				});
				this.actions.SetActionList(JSON.parse(resp.response));
			},
			err => {
				this.snackBar.open(err.error.response, "", {
					duration: 3000, panelClass: "proper-colors", horizontalPosition: 'center',
					verticalPosition: 'top',
				});
			}
		);
		form.resetForm();
	}
	public ValidateServiceName(control: AbstractControl): ValidationErrors | null {
		if (control.value != "" && control.value != undefined) {}
		return null;
	}
}
