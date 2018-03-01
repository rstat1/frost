import { Component, OnInit } from '@angular/core';

import { PageInfoService } from 'app/services/page-info.service';
import { ActionListService } from 'app/services/action-list.service';

@Component({
	selector: 'app-new-service',
	templateUrl: './new-service.html',
	styleUrls: ['./new-service.css']
})
export class NewServiceComponent implements OnInit {

	constructor(private actions: ActionListService, private header: PageInfoService) {
		
	}

	ngOnInit() {
		this.actions.OnHighlightPrimaryAction();
		this.header.SetPagePath(window.location.pathname);
	}

}
