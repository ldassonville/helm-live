import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PropertyViewerComponent } from './property-viewer.component';

describe('PropertyViewerComponent', () => {
  let component: PropertyViewerComponent;
  let fixture: ComponentFixture<PropertyViewerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PropertyViewerComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PropertyViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
