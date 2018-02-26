export class AuthRequest {
	public Username: string;
	public Password: string;
	constructor(username: string, password: string) {
		this.Username = username;
		this.Password = password;
	}
}