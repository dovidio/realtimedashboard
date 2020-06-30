import { Component, OnInit, OnDestroy } from '@angular/core';
import * as mapboxgl from 'mapbox-gl';
import { environment } from '../../environments/environment';
import { FeatureCollection } from '../model';
import { StatisticsService } from '../services/statistics.service';
import { Subscription } from 'rxjs';

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
  source: any;

  constructor(private statsService: StatisticsService) { }

  ngOnInit(): void {
    this.initializeMap();
  }

  initializeMap(): void {
    // @ts-ignore see https://github.com/DefinitelyTyped/DefinitelyTyped/issues/23467
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
        },
        cluster: true,
        clusterMaxZoom: 14, // Max zoom to cluster points on
        clusterRadius: 50 // Radius of each cluster when clustering points (defaults to 50)
      });

      this.source = this.map.getSource('appdownloads');

      this.sub = this.statsService.appDownloads$.subscribe((apps) => {
        let data = new FeatureCollection(apps);
        // @ts-ignore
        this.source.setData(data);
      });

      this.map.addLayer({
        id: 'appdownloads',
        source: 'appdownloads',
        type: 'circle',
        filter: ['has', 'point_count'],
        paint: {
          // Use step expressions (https://docs.mapbox.com/mapbox-gl-js/style-spec/#expressions-step)
          // with three steps to implement three types of circles:
          //   * Blue, 20px circles when point count is less than 100
          //   * Yellow, 30px circles when point count is between 100 and 750
          //   * Pink, 40px circles when point count is greater than or equal to 750
          'circle-color': [
            'step',
            ['get', 'point_count'],
            '#51bbd6',
            100,
            '#f1f075',
            750,
            '#f28cb1'
          ],
          'circle-radius': [
            'step',
            ['get', 'point_count'],
            20,
            100,
            30,
            750,
            40
          ]
        }
      });

      this.map.addLayer({
        id: 'appdownloads-count',
        type: 'symbol',
        source: 'appdownloads',
        filter: ['has', 'point_count'],
        layout: {
          'text-field': '{point_count_abbreviated}',
          'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
          'text-size': 12
        }
      });

      this.map.addLayer({
        id: 'unclustered-point',
        type: 'circle',
        source: 'appdownloads',
        filter: ['!', ['has', 'point_count']],
        paint: {
          'circle-color': '#11b4da',
          'circle-radius': 4,
          'circle-stroke-width': 1,
          'circle-stroke-color': '#fff'
        }
      });

    })
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }
}
