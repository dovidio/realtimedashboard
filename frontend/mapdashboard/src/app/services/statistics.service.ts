import { Injectable } from '@angular/core';
import { BehaviorSubject, timer, Subscription } from 'rxjs';
import { debounceTime, take } from 'rxjs/operators';
import { WebsocketService } from './websocket.service';
import { AppDownload, StatsSummary, GeoJson } from '../model';
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
};

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
};

const convertToGeoJson = (appDownload: AppDownload): GeoJson => {
    return new GeoJson([appDownload.longitude, appDownload.latitude], {name: appDownload.app_id, country: appDownload.country});
};

@Injectable({
  providedIn: 'root'
})
export class StatisticsService {

  private appDownloads = new BehaviorSubject<GeoJson[]>([]);
  private byCountrySubject = new BehaviorSubject<StatsSummary[]>([]);
  private byTimeOfDaySubject = new BehaviorSubject<StatsSummary[]>([]);
  private byAppSubject = new BehaviorSubject<StatsSummary[]>([]);

  appDownloads$ = this.appDownloads.asObservable().pipe(debounceTime(100));
  byCountry$ = this.byCountrySubject.asObservable();
  byTimeOfDay$ = this.byTimeOfDaySubject.asObservable();
  byApp$ = this.byAppSubject.asObservable();

  private wsSubscription = Subscription.EMPTY;

  constructor(private websocketService: WebsocketService, private appDownloadService: AppDownloadService) {
    // every five seconds
    this.setupAutomaticDataReconnection();
    this.connectAndPushDataToConsumers();
  }

  // every five seconds we check that the socket connection is healthy. If that is not the case, we try to reconnect
  // and we refresh the available data
  private setupAutomaticDataReconnection(): void {
    timer(5000, 5000).subscribe(() => {
      if (this.websocketService.socket$.isStopped || this.websocketService.socket$.closed) {
        this.websocketService.connect();
        this.wsSubscription.unsubscribe();
        this.connectAndPushDataToConsumers();
      }
    })
  }

  private connectAndPushDataToConsumers() {
    this.wsSubscription = this.websocketService.socket$.subscribe(this.updateStats.bind(this));
    this.appDownloadService.getAppDownloadList().pipe(take(1)).subscribe((appDownloads: AppDownload[]) => {
      appDownloads.forEach(this.updateStats.bind(this))
    });
  }

  private updateStats(appDownload: AppDownload) {
    const currentApps = this.appDownloads.getValue();
    currentApps.push(convertToGeoJson(appDownload));
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
