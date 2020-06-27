import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http'
import { Observable } from 'rxjs';
import { AppDownload } from '../model';

@Injectable({
  providedIn: 'root'
})
export class AppDownloadService {

  constructor(private httpclient: HttpClient) { }

  getAppDownloadList(): Observable<AppDownload[]> {
    return this.httpclient.get<AppDownload[]>("http://localhost:8080/appdownloads");
  }
}
