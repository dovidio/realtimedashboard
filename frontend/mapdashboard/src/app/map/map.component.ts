import { Component, OnInit, OnDestroy } from '@angular/core';
import mapboxgl, { GeoJSONSource } from 'mapbox-gl';
import { environment } from '../../environments/environment';
import { AppDownload, FeatureCollection } from '../model';
import { StatisticsService } from '../services/statistics.service';
import { Subscription, EMPTY } from 'rxjs';

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrls: ['./map.component.scss']
})
export class MapComponent implements OnInit, OnDestroy {
  map: mapboxgl.Map;
  style = 'mapbox://styles/umbertodov/ckbzaq1yx2muh1iqrvnhfile5';
  lat = 46.641894;
  lng = 14.269776;
  sub: Subscription = Subscription.EMPTY;
  source: mapboxgl.AnySourceImpl;
  constructor(private statsService: StatisticsService) { }

  ngOnInit(): void {
    this.initializeMap();
    this.sub = this.statsService.appDownloads$.subscribe((apps) => apps.forEach(this.addToMap.bind(this)));
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

    /// Add Marker on Click
    this.map.on('click', (event) => {
      console.log(event)
    })

    this.map.on('load', (event) => {

      this.map.addSource('appdownloads', {
        type: 'geojson',
        data: {
          type: 'FeatureCollection',
          features: []
        }
      });

      this.source = this.map.getSource('appdownloads');

      this.statsService.appDownloads$.subscribe((apps) => {
        let data = new FeatureCollection(apps);
        // @ts-ignore
        (this.source as GeoJSONSource).setData(data);
      });

      this.map.addLayer({
        id: 'appdownloads',
        source: 'appdownloads',
        type: 'symbol',
        layout: {
          'icon-image': 'marker-15',
        },
      });

    })
  }

  addToMap(appDownload: AppDownload): void {
    // new mapboxgl.Marker()
    //   .setLngLat([appDownload.longitude, appDownload.latitude])
    //   .addTo(this.map);
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }
}
