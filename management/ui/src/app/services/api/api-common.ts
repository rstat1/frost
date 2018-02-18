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