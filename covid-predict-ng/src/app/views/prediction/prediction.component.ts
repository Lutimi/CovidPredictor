import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { CovidPrediction } from 'src/app/interfaces/covid-prediction';
import { CovidPredictorService } from 'src/app/services/covid-predictor.service';

@Component({
  selector: 'app-prediction',
  templateUrl: './prediction.component.html',
  styleUrls: ['./prediction.component.css']
})
export class PredictionComponent implements OnInit {

  predictForm: FormGroup;
  region: string = ''
  nCases: number = 0;

  viewCases: boolean | undefined;
  viewPredict: CovidPrediction | undefined;

  listRegion = [
    { value: 'AMAZONAS', name: 'Amazonas' },
    { value: 'ANCASH', name: 'Áncash' },
    { value: 'APURIMAC', name: 'Apurímac' },
    { value: 'AREQUIPA', name: 'Arequipa' },
    { value: 'AYACUCHO', name: 'Ayacucho' },
    { value: 'CAJAMARCA', name: 'Cajamarca' },
    { value: 'CUSCO', name: 'Cusco' },
    { value: 'HUANCAVELICA', name: 'Huancavelica' },
    { value: 'HUANUCO', name: 'Huánuco' },
    { value: 'ICA', name: 'Ica' },
    { value: 'JUNIN', name: 'Junín' },
    { value: 'LA LIBERTAD', name: 'La Libertad' },
    { value: 'LAMBAYEQUE', name: 'Lambayeque' },
    { value: 'LIMA', name: 'Lima' },
    { value: 'LORETO', name: 'Loreto' },
    { value: 'MADRE DE DIOS', name: 'Madre de Dios' },
    { value: 'MOQUEGUA', name: 'Moquegua' },
    { value: 'PASCO', name: 'Pasco' },
    { value: 'PIURA', name: 'Piura' },
    { value: 'PUNO', name: 'Puno' },
    { value: 'SAN MARTIN', name: 'San Martín' },
    { value: 'TACNA', name: 'Tacna' },
    { value: 'TUMBES', name: 'Tumbes' },
    { value: 'UCAYALI', name: 'Ucayali' }
  ];

  selectRegion = this.listRegion[0].name;

  constructor(private covidPredictionService: CovidPredictorService, private dialogRef: MatDialogRef<PredictionComponent>, 
    private formBuilder: FormBuilder) { 
      this.predictForm = this.formBuilder.group({
        dateAt: ['YYYYMMDD', [Validators.required]]
      })
    }

  ngOnInit(): void { 
    this.viewCases = false;
  }

  predictCovid() {
    // console.log(this.convertString(this.predictForm.value.dateAt))
    this.covidPredictionService.predictCovidCase(this.selectRegion, this.convertString(this.predictForm.value.dateAt)).subscribe({
      error: (err) => console.log(err),
      next: (rest) => {
        this.viewPredict = rest;
        console.log(this.viewPredict);
        this.viewCases = true;
      },
      complete: () => console.log('Complete')
    });
  }

  addPredict() {
    console.log(this.viewPredict)
    this.covidPredictionService.savePrediction(this.viewPredict).subscribe({
      error: (err) => console.log(err),
      next: (rest) => {
        console.log('Successfull Save Predict')
        this.dialogRef.close();
        window.location.href = '/covidPredict'
      },
      complete: () => console.log('Complete')
    });
  }

  convertString(dateName: string) {
    var year = dateName.substr(0, 4)
    var month = dateName.substr(5, 2)
    var day = dateName.substr(8, 2)

    var newDate = new String(year + month + day)
    return newDate
  }

}
