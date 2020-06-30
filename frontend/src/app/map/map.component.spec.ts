import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MapComponent } from './map.component';
import { StatisticsService } from '../services/statistics.service';
import { MockedStatisticsService } from '../mocks/mocks.spec';
import { EMPTY } from 'rxjs';

describe('MapComponent', () => {
  let component: MapComponent;
  let fixture: ComponentFixture<MapComponent>;


  const mockedStatsService = new MockedStatisticsService(EMPTY, EMPTY, EMPTY, EMPTY);

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MapComponent ],
      providers: [
        { provide: StatisticsService, useValue: mockedStatsService },
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MapComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
