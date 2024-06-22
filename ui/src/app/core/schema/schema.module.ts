import { NgModule } from '@angular/core';
import { SchemaViewerComponent } from './schema-viewer/schema-viewer.component';
import { RefViewerComponent } from './ref-viewer/ref-viewer.component';
import { JsonPipe, KeyValuePipe, NgForOf } from '@angular/common';
import { PropertyViewerComponent } from './property-viewer/property-viewer.component';

@NgModule({
  exports: [
    SchemaViewerComponent,
    RefViewerComponent,
    PropertyViewerComponent
  ],
  declarations: [
    SchemaViewerComponent,
    RefViewerComponent,
    PropertyViewerComponent,
  ],
  imports: [
    JsonPipe, KeyValuePipe, NgForOf
  ],
})
export class SchemaModule {}
