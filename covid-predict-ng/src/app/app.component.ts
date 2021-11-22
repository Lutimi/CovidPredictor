import { Component, OnInit, ViewChild } from '@angular/core';
import { CovidPredictorService } from './services/covid-predictor.service';
import { CovidPrediction } from './interfaces/covid-prediction';
import { MatDialog } from '@angular/material/dialog';
import { PredictionComponent } from './views/prediction/prediction.component';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {

  covidPredicts: CovidPrediction[] = [];
  displayItems: string[] = ["region", "nCases", "date"];

  constructor(private covidPredictionService: CovidPredictorService, private dialog: MatDialog) {}

  ngOnInit(): void {
    this.viewAllPredictis();
    // this.covidPredicts.sort = this.sort;
    // this.covidPredicts.paginator = this.paginator;
  }

  viewAllPredictis() {
    this.covidPredictionService.listPredictions().subscribe({
      error: (err) => console.log(err),
      next: (rest) => {
        this.covidPredicts = rest;
        console.log(this.covidPredicts);
      },
      complete: () => console.log('Complete')
    })
  }

  title = 'covid-predict-ng';

  convertDate(dateName: string) {
    var year = dateName.substr(0, 4)
    var month = dateName.substr(4, 2)
    var day = dateName.substr(6, 2)

    var newDate = new Date(year + "-" + month + "-" + day)
    return newDate.toLocaleDateString("es-PE")
  }

  generatePrediction() {
    this.dialog.open(PredictionComponent, { autoFocus: true })
  }
}
