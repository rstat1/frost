export interface APIResponse {
	status: string;
	response: string;
}
export interface Service {
	name: string;
	filename: string;
	api_prefix: string;
	address: string;
	managed: boolean;
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
	permissions: ServiceAccess[];
	constructor(name: string, password: string, permissions: ServiceAccess[])
	{
		this.username = name;
		this.password = password;
		this.permissions = permissions;
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