import { Subject ,  Observable } from 'rxjs';
import { Injectable } from '@angular/core';

@Injectable()
export class PageInfoService {
	public PagePath: Observable<string>;
	public PageTitle: Observable<string>;
	public PageLogo: Observable<string>;

	private pagePath: Subject<string>;
	private pageTitle: Subject<string>;
	private pageLogo: Subject<string>;

	constructor() {
		this.pagePath = new Subject<string>();
		this.pageTitle = new Subject<string>();
		this.pageLogo = new Subject<string>();

		this.PagePath = this.pagePath.asObservable();
		this.PageTitle = this.pageTitle.asObservable();
		this.PageLogo = this.pageLogo.asObservable();

	}
	public SetPagePath(newPath: string) {
		this.pagePath.next(newPath);
	}
	public SetPageLogoAndTitle(newLogoName: string, newTitle: string) {
		this.pageLogo.next(newLogoName);
		this.pageTitle.next(newTitle);
	}
}