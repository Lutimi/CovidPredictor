import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { CovidPrediction } from '../interfaces/covid-prediction';

@Injectable({
  providedIn: 'root'
})
export class CovidPredictorService {

  constructor(private http: HttpClient) { }

  // Use Algorithm
  predictCovidCase(region: string, date: String): Observable<CovidPrediction> {
    let params = new HttpParams();
    params = params.set('region', region)
    params = params.set('date', String(date))
    return this.http.get<CovidPrediction>(environment.apiPredict + 'make-prediction', {params});
  }

  // List All Predictions
  listPredictions(): Observable<CovidPrediction[]> {
    return this.http.get<CovidPrediction[]>(environment.apiMicroservices + 'list-predictions');
  }

  // Save Prediction
  savePrediction(body: CovidPrediction | undefined): Observable<CovidPrediction>{
    return this.http.post<CovidPrediction>(environment.apiMicroservices + 'save-prediction', body);
  }
}
