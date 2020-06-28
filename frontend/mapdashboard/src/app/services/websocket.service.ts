import { Injectable } from '@angular/core';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { environment } from 'src/environments/environment';
import { AppDownload } from '../model';

const WS_ENDPOINT = environment.wsEndpoint;

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  public socket$: WebSocketSubject<AppDownload>;

  constructor() { 
    this.connect();
  }

  public connect(): void {
    if (!this.socket$ || this.socket$.closed) {
      this.socket$ = webSocket(WS_ENDPOINT);
    }
  }

  close(): void {
    !this.socket$ || this.socket$.complete();
  }
}
