export class Project {
	projectName: string;
	hasSubModules: boolean;
	subModules: string[];
	latestCommit: Commit;
}
export class Commit {
	author: string;
	branch: string;
	hash: string;
	message: string;
	time: string;
}
export interface ProjectsQueryResponse {
	projects: Project[];
}
export interface ProjectQueryResponse {
	project: Project;
}
export interface TopCommitResponse {
	latestCommit: Commit;
}
export interface AllCommitsResponse {
	commits: Commit[];
}