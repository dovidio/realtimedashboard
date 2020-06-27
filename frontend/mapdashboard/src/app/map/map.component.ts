import { Component, OnInit } from '@angular/core';
import mapboxgl from 'mapbox-gl';
import { environment } from '../../environments/environment';
import { WebsocketService } from '../services/websocket.service';
import { AppDownloadService } from '../services/app-download.service';
import { AppDownload } from '../model';
import { tap } from 'rxjs/operators';

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

  constructor(private appDownloadService: AppDownloadService, private websocketService: WebsocketService) { }

  ngOnInit(): void {
    this.initializeMap();
    this.startWebsocket();
    this.makeFirstRequest();
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

  startWebsocket(): void {
    this.websocketService.connect();
    this.websocketService.socket$
      .pipe(tap((e) => console.log("called but is there?", e)))
      .subscribe(this.addToMap.bind(this));
  }

  makeFirstRequest(): void {
    this.appDownloadService.getAppDownloadList().subscribe((appDownloads: AppDownload[]) => {
      appDownloads.forEach(this.addToMap.bind(this))
    });
  }

  addToMap(appDownload: AppDownload): void {
    new mapboxgl.Marker()
      .setLngLat([appDownload.longitude, appDownload.latitude])
      .addTo(this.map);
  }
}
