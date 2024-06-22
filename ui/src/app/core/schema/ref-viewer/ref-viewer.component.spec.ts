import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RefViewerComponent } from './ref-viewer.component';

describe('RefViewerComponent', () => {
  let component: RefViewerComponent;
  let fixture: ComponentFixture<RefViewerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RefViewerComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(RefViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
