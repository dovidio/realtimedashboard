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