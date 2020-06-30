import { Component } from '@angular/core';
import { StatisticsService } from '../services/statistics.service';


@Component({
  selector: 'app-statistics',
  templateUrl: './statistics.component.html',
})
export class StatisticsComponent {

  selectedGrouping: "bycountry" | "bytimeofday" | "byapp" = "bycountry";

  constructor(private statsService: StatisticsService) { }

  byCountry$ = this.statsService.byCountry$;
  byTimeOfDay$ = this.statsService.byTimeOfDay$;
  byApp$ = this.statsService.byApp$;
}
