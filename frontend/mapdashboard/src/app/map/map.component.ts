import { Component, OnInit } from '@angular/core';
import mapboxgl from 'mapbox-gl';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { AppDownloadService } from '../services/app-download.service';
import { AppDownload } from '../model';


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

  constructor(private appDownloadService: AppDownloadService) { }

  ngOnInit(): void {
    mapboxgl.accessToken = environment.mapbox.accessToken;
    this.map = new mapboxgl.Map({
      container: 'map',
      style: this.style,
      zoom: 4,
      center: [this.lng, this.lat]
    });
    // Add map controls
    this.map.addControl(new mapboxgl.NavigationControl());

    this.appDownloadService.getAppDownloadList().subscribe((appDownloads: AppDownload[]) => {

      appDownloads.forEach((appDownload) => {
        var marker = new mapboxgl.Marker()
          .setLngLat([appDownload.longitude, appDownload.latitude])
          .addTo(this.map);
      })
    });
  }

}
