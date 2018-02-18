import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { ConfigService } from "app/services/config.service";
import { APIResponse } from 'app/services/api/api-common';

class GraphQLQuery {
	variables: any;
	query: string;
	operationName: string;
}
export class GraphQLObject {
	objectName: string;
	fields?: string[];
}

@Injectable()
export class GraphQLService {
	constructor(private http: HttpClient) {}

	public Query<T>(fieldName: string, args: Map<string, any>, fields: GraphQLObject[]): Observable<T> {
		var subject: Subject<T> = new Subject<T>();
		var apiURL: string = ConfigService.GetAPIURLFor("query");
		var query: string = this.MakeQueryWithObject(fieldName, args, fields, false);
		this.http.post<APIResponse>(apiURL, query).subscribe(data => {
			if (data.status == "success") {
				let resp = JSON.parse(data.response)
				subject.next(resp.data as T)
			}
		}, error => { subject.next(null); })
		return subject
	}
	public Mutation(fieldName: string, args: Map<string, any>, fields: GraphQLObject[]): Observable<boolean> {
		var subject: Subject<boolean> = new Subject<boolean>();
		var apiURL: string = ConfigService.GetAPIURLFor("query");
		var query: string = this.MakeQueryWithObject(fieldName, args, fields, true);
		this.http.post<APIResponse>(apiURL, query).subscribe(data => {
			if (data.status == "success") {
				subject.next(true);
			} else if (data.status == "failed") {
				subject.next(false);
			}
		})
		return subject
	}
	private MakeGQLQueryString(query: string): string {
		var gqlQuery: GraphQLQuery = {
			variables: null,
			query: query,
			operationName: "",
		};
		return JSON.stringify(gqlQuery);
	}
	private MakeQuery(name: string, args: Map<string, any>, elements: string[]): string {
		var index: number = 0;
		var query: string = "{ " + name ;
		if (args != null) {
			query += "("
			args.forEach((value, key, m) => {
				query += key + ": " + value;
				if (index < args.size - 1)
				{
					query += ", ";
				}
				index++;
			})
			query += ")"
		}
		query += " { "
		elements.forEach(element => {
			query += element + " "
		});
		query += "} }"
		return this.MakeGQLQueryString(query);
	}
	private MakeQueryWithObject(name: string, args: Map<string, any>, fields: GraphQLObject[], isMutation: boolean): string {
		var index: number = 0;
		var elementHasFields: boolean = false;
		var query: string = "{ " + name ;
		if (isMutation) { query  = "mutation { " + name ; }
		if (args != null) {
			query += "("
			args.forEach((value, key, m) => {
				query += key + ": " + value;
				if (index < args.size - 1)
				{
					query += ", ";
				}
				index++;
			})
			query += ")"
		}
		query += " { "
		fields.forEach(element => {
			query += element.objectName;
			if (element.fields != null) {
				elementHasFields = true;
				query += " { ";
				element.fields.forEach(element => {
					query += element + " "
				});
			} else {
				query += " "
			}
		});
		if (elementHasFields) { query += "} } }"; }
		else { query += "} }"; }
		return this.MakeGQLQueryString(query);
	}
}