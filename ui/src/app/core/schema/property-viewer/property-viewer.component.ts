import { Component, Input, OnInit, input } from '@angular/core';

@Component({
  selector: 'app-property-viewer',
  templateUrl: './property-viewer.component.html',
  styleUrl: './property-viewer.component.scss'
})
export class PropertyViewerComponent  implements OnInit {

  @Input() mode: 'inline'|'standalone' = 'standalone'

  @Input() property: any;

  @Input() root: boolean = false

  @Input() name: string = ""

  @Input() file: string = ""

  @Input() required: boolean = false

  @Input() collapsed: boolean = true  ;

  constructor(){}

  ngOnInit(): void {
    if(this.mode == 'inline'){
      this.collapsed = false
    }
  } 

  onClick() {
    this.collapsed = !this.collapsed;
  }

  isRequired(name: string): boolean {
    if (!this.property){
      return false
    }

    if (this.property?.required && this.property?.required.length > 0) {
      return this.property.required.includes(name)
    }
    return false
  }

  
}

