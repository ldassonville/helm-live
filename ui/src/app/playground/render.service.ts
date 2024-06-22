import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {HelmRender} from "./render";

@Injectable({
  providedIn: 'root'
})
export class RenderService {

    constructor(private http: HttpClient) { }

    baseUrl = "http://localhost:8085"
    render(): Observable<HelmRender> {
      return this.http.get<HelmRender>(`${this.baseUrl}/_render`)
    }

  live(): Observable<HelmRender> {

    const eventSource = new EventSource(`${this.baseUrl}/stream`);

    return new Observable(observer => {
      eventSource.addEventListener('message', event => {
        console.log("RenderService: live: addEventListener: message: ", event.data);
        const messageData: HelmRender = JSON.parse(event.data);
        observer.next(messageData);

      })

    });

  }
}
