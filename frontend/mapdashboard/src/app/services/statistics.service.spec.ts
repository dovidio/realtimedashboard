import { TestBed } from '@angular/core/testing';
import { EMPTY, of, Subject } from 'rxjs';
import { mockedAppDownloads, MockedAppDownloadService, MockedWebsocketService } from '../mocks/mocks.spec';
import { AppDownload } from '../model';
import { AppDownloadService } from './app-download.service';
import { StatisticsService } from './statistics.service';
import { WebsocketService } from './websocket.service';

describe('StatisticsService', () => {
  let service: StatisticsService;
  let mockWebSocketStream;

  beforeEach(() => {
    mockWebSocketStream = new Subject<AppDownload>();

    const mockedAppDownloadService = new MockedAppDownloadService(of(mockedAppDownloads));
    const mockedWebsocketService = new MockedWebsocketService(mockWebSocketStream);
    

    TestBed.configureTestingModule({
      providers: [
        { provide: WebsocketService, useValue: mockedWebsocketService },
        { provide: AppDownloadService, useValue: mockedAppDownloadService },
      ],
    });
    service = TestBed.inject(StatisticsService);
  });

  it('should load the initial data correctly from the appdownloadservice', () => {

    service.appDownloads$.subscribe((apps) => {
      expect(apps.length).toBe(mockedAppDownloads.length);
    });

    service.byCountry$.subscribe((byCountry) => {
      expect(byCountry.find((c) => c.name == "Austria").count).toBe(2);
      expect(byCountry.find((c) => c.name == "USA").count).toBe(2);
      expect(byCountry.find((c) => c.name == "Italy").count).toBe(1);
    });

    service.byTimeOfDay$.subscribe((byTimeOfDay) => {
      expect(byTimeOfDay.find((c) => c.name == "morning").count).toBe(1);
      expect(byTimeOfDay.find((c) => c.name == "afternoon").count).toBe(1);
      expect(byTimeOfDay.find((c) => c.name == "evening").count).toBe(2);
      expect(byTimeOfDay.find((c) => c.name == "night").count).toBe(1);
    });

    service.byApp$.subscribe((byApp) => {
      expect(byApp.find((c) => c.name == "IOS_ALERT").count).toBe(4);
      expect(byApp.find((c) => c.name == "ANDROID_ALERT").count).toBe(1);
    });
  });

  it('should update data when receing update from the websocket', () => {
    // given
    const update = {
      latitude: 0,
      longitude: 0,
      country: "Italy",
      app_id: "ANDROID_ALERT",
      downloaded_at: 1593318274749
    }

    // when
    mockWebSocketStream.next(update);

    // then
    service.appDownloads$.subscribe((apps) => {
      expect(apps.length).toBe(mockedAppDownloads.length + 1);
    });

    service.byCountry$.subscribe((byCountry) => {
      expect(byCountry.find((c) => c.name == "Austria").count).toBe(2);
      expect(byCountry.find((c) => c.name == "USA").count).toBe(2);
      expect(byCountry.find((c) => c.name == "Italy").count).toBe(1 + 1);
    });

    service.byTimeOfDay$.subscribe((byTimeOfDay) => {
      expect(byTimeOfDay.find((c) => c.name == "morning").count).toBe(1);
      expect(byTimeOfDay.find((c) => c.name == "afternoon").count).toBe(1);
      expect(byTimeOfDay.find((c) => c.name == "evening").count).toBe(2);
      expect(byTimeOfDay.find((c) => c.name == "night").count).toBe(1 + 1);
    });

    service.byApp$.subscribe((byApp) => {
      expect(byApp.find((c) => c.name == "IOS_ALERT").count).toBe(4);
      expect(byApp.find((c) => c.name == "ANDROID_ALERT").count).toBe(1 + 1);
    });
  })
});
