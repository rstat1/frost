import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';


@Injectable()
export class ConfigService {
	private static API_VERSION_TAG: string = "trinity";
	private static API_ENDPOINT: string = environment.APIBaseURL;
	private static FROST_ENDPOINT: string = environment.FrostBaseURL;

	public static GetAPIURLFor(api: string): string {
		return this.API_ENDPOINT + "/" + this.API_VERSION_TAG + "/" + api;
	}
	public static GetFrostURLFor(api: string): string {
		return this.FROST_ENDPOINT + "/" + api;
	}
}