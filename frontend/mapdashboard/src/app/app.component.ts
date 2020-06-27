import { Component, OnInit } from "@angular/core";
import { AppDownloadService } from './services/app-download.service';

@Component({
    templateUrl: "app.component.html",
    selector: 'app-root',
    styleUrls: ["app.component.scss"]
})
export class AppComponent implements OnInit {

    constructor(private appDownloadService: AppDownloadService) {}

    ngOnInit(): void {
        this.appDownloadService.getAppDownloadList().subscribe(console.log);
    }
}
