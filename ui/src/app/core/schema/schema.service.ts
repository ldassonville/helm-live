import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SchemaService {

  constructor(private http: HttpClient) { }

  getSchema(schemaUrl: string): Observable<any> {

    return this.http.get(schemaUrl, {
      responseType: "json"   // This one here tells HttpClient to parse it as text, not as JSON
    });
  }
}

