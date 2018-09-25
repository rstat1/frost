import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { APIService } from 'app/services/api/api.service';
import { RouteAlias } from 'app/services/api/api-common';

@Component({
	selector: 'app-new-alias-dialog',
	templateUrl: './new-alias.html',
	styleUrls: ['./new-alias.css']
})
export class NewAliasDialogComponent {
	public forAPI: string;
	public aliasURL: string = "";
	public aliasedRoute: string = "";

	constructor(@Inject(MAT_DIALOG_DATA) public data: any, public dialogRef: MatDialogRef<NewAliasDialogComponent>,
		private api: APIService) {
		this.forAPI = <string>data.apiName;
	}
	public save() {
		let routeAlias: RouteAlias = new RouteAlias();
		routeAlias.apiName = this.forAPI;
		routeAlias.fullURL = this.aliasURL;
		routeAlias.apiRoute = this.aliasedRoute;
		this.api.NewRouteAlias(routeAlias).subscribe(
			resp => {
				if (resp.status == "success") { this.dialogRef.close(resp); }
			},
			e => { this.dialogRef.close(e.error); }
		);
	}
}