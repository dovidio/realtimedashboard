export interface AppDownload {
    latitude: number;
    longitude: number;
    app_id: string;
    downloaded_at: number;
    country: string;
}

export interface StatsSummary {
    name: string;
    count: number;
}

export interface Geometry {
    type: string;
    coordinates: number[];
}

export interface IGeoJson {
    type: string;
    geometry: Geometry;
    properties?: any;
    $key?: string;
}

export class GeoJson implements IGeoJson {
  type = 'Feature';
  geometry: Geometry;

  constructor(coordinates, public properties?) {
    this.geometry = {
      type: 'Point',
      coordinates: coordinates
    }
  }
}

export class FeatureCollection {
  type = 'FeatureCollection'
  constructor(public features: Array<GeoJson>) {}
}