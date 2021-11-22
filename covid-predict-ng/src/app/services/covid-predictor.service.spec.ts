import { TestBed } from '@angular/core/testing';

import { CovidPredictorService } from './covid-predictor.service';

describe('CovidPredictorService', () => {
  let service: CovidPredictorService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CovidPredictorService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
