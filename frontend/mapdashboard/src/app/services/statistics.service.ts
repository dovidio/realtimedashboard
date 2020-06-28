import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { debounceTime } from 'rxjs/operators';
import { WebsocketService } from './websocket.service';
import { AppDownload, StatsSummary } from '../model';
import { AppDownloadService } from './app-download.service';

type TimeOfDay = "morning" | "afternoon" | "evening" | "night";

const incrementStats = (array: StatsSummary[], prop: string): void => {
  const index = array.findIndex((v) => v.name === prop);
  if (index < 0) {
    array.push({name: prop, count: 1});
  } else {
    array[index].count++;
  }
  array.sort((a, b) => b.count - a.count);
}

const getTimeOfDay = (appDownload: AppDownload): TimeOfDay => {

  const hours = new Date(appDownload.downloaded_at).getUTCHours();

  if (hours > 5 && hours <= 12) {
    return "morning";
  }
  if (hours > 12 && hours <= 17) {
    return "afternoon";
  }
  if (hours > 17 && hours <= 23) {
    return "evening";
  }

  return "night";
}

@Injectable({
  providedIn: 'root'
})
export class StatisticsService {

  private appDownloads = new BehaviorSubject<AppDownload[]>([]);
  private byCountrySubject = new BehaviorSubject<StatsSummary[]>([]);
  private byTimeOfDaySubject = new BehaviorSubject<StatsSummary[]>([]);
  private byAppSubject = new BehaviorSubject<StatsSummary[]>([]);

  appDownloads$ = this.appDownloads.asObservable().pipe(debounceTime(200));
  byCountry$ = this.byCountrySubject.asObservable();
  byTimeOfDay$ = this.byTimeOfDaySubject.asObservable();
  byApp$ = this.byAppSubject.asObservable();

  constructor(private websocketService: WebsocketService, private appDownloadService: AppDownloadService) {
    this.websocketService.socket$.pipe(debounceTime(10)).subscribe(this.updateStats.bind(this));
    this.appDownloadService.getAppDownloadList().subscribe((appDownloads: AppDownload[]) => {
      appDownloads.forEach(this.updateStats.bind(this))
    });
  }

  private updateStats(appDownload: AppDownload) {
    const currentApps = this.appDownloads.getValue();
    currentApps.push(appDownload);
    this.appDownloads.next(currentApps);
    this.updateCountryStatistics(appDownload);
    this.updateTimeOfDaySatistics(appDownload);
    this.updateAppStatistics(appDownload);
  }

  private updateCountryStatistics(appDownload: AppDownload): void {
    const next = this.byCountrySubject.getValue();
    incrementStats(next, appDownload.country);
    this.byCountrySubject.next(next);
  }

  private updateTimeOfDaySatistics(appDownload: AppDownload): void {
    const timeOfDay = getTimeOfDay(appDownload);
    const next = this.byTimeOfDaySubject.getValue();
    incrementStats(next, timeOfDay);
    this.byTimeOfDaySubject.next(next);
  }

  private updateAppStatistics(appDownload: AppDownload): void {
    const app = appDownload.app_id;
    const next = this.byAppSubject.getValue();
    incrementStats(next, app);
    this.byAppSubject.next(next);
  }
}
