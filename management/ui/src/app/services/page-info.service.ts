import { Subject } from 'rxjs/Subject';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class PageInfoService {
	public PagePath: Observable<string>;
	public PageTitle: Observable<string>;

	private pagePath: Subject<string>;
	private pageTitle: Subject<string>;

	constructor() {
		this.pagePath = new Subject<string>();
		this.pageTitle = new Subject<string>();

		this.PagePath = this.pagePath.asObservable();
		this.PageTitle = this.pageTitle.asObservable();
	}
	public SetPagePath(newPath: string) {
		this.pagePath.next(newPath);
	}
	public SetPageTitle(newTitle: string) {
		this.pageTitle.next(newTitle);
	}
}
