import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { By } from '@angular/platform-browser';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { BehaviorSubject, EMPTY } from 'rxjs';
import { mockedByAppStats, mockedByCountryStats, mockedByTimeOfDayStats, MockedStatisticsService } from '../mocks/mocks.spec';
import { StatsSummary } from '../model';
import { StatisticsService } from '../services/statistics.service';
import { StatisticsComponent } from './statistics.component';

describe('StatisticsComponent', () => {
  let component: StatisticsComponent;
  let fixture: ComponentFixture<StatisticsComponent>;

  const byCountrySubject = new BehaviorSubject<StatsSummary[]>(mockedByCountryStats);
  const byTimeOfDaySubject = new BehaviorSubject<StatsSummary[]>(mockedByTimeOfDayStats);
  const byAppSubject = new BehaviorSubject<StatsSummary[]>(mockedByAppStats);

  const mockedStatsService = new MockedStatisticsService(EMPTY, byCountrySubject, byTimeOfDaySubject, byAppSubject);

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [StatisticsComponent],
      imports: [
        MatSelectModule,
        MatFormFieldModule,
        NoopAnimationsModule,
      ],
      providers: [
        { provide: StatisticsService, useValue: mockedStatsService }
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(StatisticsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should show statistics by country when selected', () => {
    // when
    component.selectedGrouping = 'bycountry';
    fixture.detectChanges();

    // then
    const countryList = fixture.debugElement.query(By.css('[uitestid="countryStatsList"]'));
    expect(countryList).toBeTruthy();
  });

  it('should show statistics by time of day when selected', () => {
    // when
    component.selectedGrouping = 'bytimeofday';
    fixture.detectChanges();

    // then
    const timeOfDayStatsList = fixture.debugElement.query(By.css('[uitestid="timeOfDayStatsList"]'));
    expect(timeOfDayStatsList).toBeTruthy();
  });

  it('should show statistics by app when selected', () => {
    // when
    component.selectedGrouping = 'byapp';
    fixture.detectChanges();

    // then
    const appStatsList = fixture.debugElement.query(By.css('[uitestid="appStatsList"]'));
    expect(appStatsList).toBeTruthy();
  });
});
