export interface APIResponse {
	status: string;
	response: string;
}
export class Service {
	public name: string;
	public address: string;
	public filename: string;
	public managed: boolean;
	public api_prefix: string;
	public RedirectURL: string;
	public ServiceKey: string;
	public ServiceID: string;
	public managedUpdates: boolean;
	constructor() {
		this.name = "";
		this.filename = "";
	}
}
export interface ServiceAccess {
	service: string;
	permissions: Permission[];
}
export interface Permission {
	name: string;
	value: boolean;
}
export class NewUser {
	username: string;
	password: string;
	access: ServiceAccess[];
	constructor(name: string, password: string, permissions: ServiceAccess[]) {
		this.username = name;
		this.password = password;
		this.access = permissions;
	}
}
export class AuthRequest {
	public Username: string;
	public Password: string;
	constructor(username: string, password: string) {
		this.Username = username;
		this.Password = password;
	}
}
export interface PermissionValues {
	p2: boolean;
	p3: boolean;
}

export interface Service {
	p: PermissionValues;
}
export class ServiceEdit {
	public name: string;
	public property: string;
	public new: string;
}
export class RouteAlias {
	public apiName: string;
	public fullURL: string;
	public apiRoute: string;
}