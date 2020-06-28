import { Component, OnInit } from '@angular/core';
import mapboxgl from 'mapbox-gl';
import { environment } from '../../environments/environment';
import { AppDownload } from '../model';
import { StatisticsService } from '../services/statistics.service';

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrls: ['./map.component.scss']
})
export class MapComponent implements OnInit {
  map: mapboxgl.Map;
  style = 'mapbox://styles/mapbox/streets-v11';
  lat = 46.641894;
  lng = 14.269776;

  constructor(private statsService: StatisticsService) { }

  ngOnInit(): void {
    this.initializeMap();
    this.statsService.appDownloads$.subscribe((apps) => apps.forEach(this.addToMap.bind(this)));
  }

  initializeMap(): void {
    mapboxgl.accessToken = environment.mapbox.accessToken;
    this.map = new mapboxgl.Map({
      container: 'map',
      style: this.style,
      zoom: 4,
      center: [this.lng, this.lat]
    });
    // Add map controls
    this.map.addControl(new mapboxgl.NavigationControl());
    this.map.addControl(new mapboxgl.ScaleControl())
  }

  addToMap(appDownload: AppDownload): void {
    new mapboxgl.Marker()
      .setLngLat([appDownload.longitude, appDownload.latitude])
      .addTo(this.map);
  }
}
