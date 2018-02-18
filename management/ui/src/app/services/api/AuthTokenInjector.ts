import { Injectable } from '@angular/core';
import { HttpEvent, HttpInterceptor, HttpHandler, HttpRequest } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';

import { ConfigService } from "app/services/config.service";

@Injectable()
export class AuthTokenInjector implements HttpInterceptor {
	intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
		req = req.clone({ setHeaders: {	Authorization: `Bearer ${ConfigService.GetAccessToken()}` } })
		return next.handle(req);
  	}
}
