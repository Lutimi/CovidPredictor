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
  predictCovidCase(region: string, date: string): Observable<CovidPrediction> {
    let params = new HttpParams();
    params = params.set('region', region)
    params = params.set('date', date)
    return this.http.get<CovidPrediction>(environment.apiURL + '1/make-prediction', {params});
  }

  // List All Predictions
  listPredictions(): Observable<CovidPrediction[]> {
    return this.http.get<CovidPrediction[]>(environment.apiURL + '0/list-predictions');
  }

  // Save Prediction
  savePrediction(body: CovidPrediction): Observable<CovidPrediction>{
    return this.http.post<CovidPrediction>(environment.apiURL + '0/save-prediction', body);
  }
}
