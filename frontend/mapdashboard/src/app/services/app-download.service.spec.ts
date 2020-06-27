import { TestBed } from '@angular/core/testing';

import { AppDownloadService } from './app-download.service';

describe('AppDownloadService', () => {
  let service: AppDownloadService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AppDownloadService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
