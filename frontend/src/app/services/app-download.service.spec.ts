import { TestBed } from '@angular/core/testing';

import { AppDownloadService } from './app-download.service';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { environment } from 'src/environments/environment';

describe('AppDownloadService', () => {
  let service: AppDownloadService;
  let controller: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
    });
    service = TestBed.inject(AppDownloadService);

    controller = TestBed.get(HttpTestingController);
  });

  it('should call the correct endpoint', () => {
    // when
    service.getAppDownloadList().subscribe();

    // then
    controller.expectOne(`${environment.baseEndpoint}/appdownloads`);
  });
});
