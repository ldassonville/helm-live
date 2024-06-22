import {Component, Input, OnInit} from '@angular/core';
import {SchemaResolver} from "../../core/schema/schema-resolver";
import {SchemaModule} from "../../core/schema/schema.module";

@Component({
  selector: 'app-schema-viewer',
  standalone: true,
  imports: [
    SchemaModule
  ],
  templateUrl: './schema-viewer.component.html',
  styleUrl: './schema-viewer.component.scss'
})
export class SchemaViewerComponent implements OnInit{

  @Input()
  public gvk: {
    group: string,
    kind: string,
    version: string
  }| undefined


  protected schema: any
  protected schemaUrl: string = ""

  constructor(private schemaResolver: SchemaResolver) {
  }

  ngOnInit(): void {
    this.loadSchema()
  }



  loadSchema(){
    this.schema = undefined
    this.schemaResolver
      .init(this.schemaUrl)
      .subscribe(schema => {
          this.schema = schema;
        }
      )
  }
}
