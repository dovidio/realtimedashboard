import { Component, OnInit } from "@angular/core";
import { AppDownloadService } from './services/app-download.service';

@Component({
    template: `<h1>Hello World!</h1>`,
    selector: 'app-root'
})
export class AppComponent implements OnInit {

    constructor(private appDownloadService: AppDownloadService) {}

    ngOnInit(): void {
        this.appDownloadService.getAppDownloadList().subscribe(console.log);
    }
}
