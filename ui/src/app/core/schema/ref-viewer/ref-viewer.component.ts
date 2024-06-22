import { Component, Input, OnInit } from '@angular/core';
import { SchemaRef, SchemaResolver } from '../schema-resolver';

@Component({
  selector: 'app-ref-viewer',
  templateUrl: './ref-viewer.component.html',
  styleUrl: './ref-viewer.component.scss'
})
export class RefViewerComponent implements OnInit {

  @Input() file: string = ""

  @Input() ref: string = ""

  refFile: string = ""

  absoluteUrl: string = ""

  schema: any
  target: any

  property: any

  constructor(private resolver: SchemaResolver){}

  ngOnInit(): void {
    this.resolveRef(this.ref)
  }

  resolveRef(refStr: string) {

    console.log(refStr)

    this.absoluteUrl = this.resolver.getAbsoluteURL(this.file, this.ref)

    this.refFile = new SchemaRef(refStr, this.file).schemaPath.filepath

    this.resolver
      .resolveRef(this.file, refStr)
      .subscribe(r => {
        this.schema = r.schema
        this.target = r.target
      })
  }
}
