import { Observable } from 'rxjs';
import { GeoJson, StatsSummary, AppDownload } from '../model';

export class MockedStatisticsService {
    constructor(
        public appDownloads$: Observable<GeoJson[]>,
        public byCountry$: Observable<StatsSummary[]>,
        public byTimeOfDay$: Observable<StatsSummary[]>,
        public byApp$: Observable<StatsSummary[]>) {}
}

export class MockedAppDownloadService {
    constructor(private appDownloads$: Observable<AppDownload[]>) {}

    getAppDownloadList(): Observable<AppDownload[]> {
        return this.appDownloads$;
    }
}

export class MockedWebsocketService {
    constructor(public socket$: Observable<AppDownload>) {}
    connect(): void {}
    disconnect(): void {}
}

export const mockedGeoJson: GeoJson[] = [
    new GeoJson([0, 0]),
    new GeoJson([1, 0]),
    new GeoJson([0, 1]),
];

export const mockedByCountryStats: StatsSummary[] = [
    {name: "Austria", count: 10},
    {name: "Italy", count: 20},
];

export const mockedByTimeOfDayStats: StatsSummary[] = [
    {name: "morning", count: 10},
    {name: "afternoon", count: 20},
];

export const mockedByAppStats: StatsSummary[] = [
    {name: "IOS_ALERT", count: 10},
    {name: "ANDROID_ALERT", count: 20},
];

export const mockedAppDownloads: AppDownload[] = [
    {
        latitude: 0,
        longitude: 0,
        country: "Austria",
        app_id: "IOS_ALERT",
        downloaded_at: 1593368544952
    },
    {
        latitude: 0,
        longitude: 0,
        country: "Austria",
        app_id: "IOS_ALERT",
        downloaded_at: 1593339845727
    },
    {
        latitude: 0,
        longitude: 0,
        country: "Italy",
        app_id: "IOS_ALERT",
        downloaded_at: 1593318274749
    },
    {
        latitude: 0,
        longitude: 0,
        country: "USA",
        app_id: "IOS_ALERT",
        downloaded_at: 1593356249476
    },
    {
        latitude: 0,
        longitude: 0,
        country: "USA",
        app_id: "ANDROID_ALERT",
        downloaded_at: 1593368544952
    },
];