import { Component, ElementRef, EventEmitter, HostListener, Input, Output } from '@angular/core';

@Component({
  selector: 'app-sliding-panel',
  templateUrl: './sliding-panel.component.html',
  styleUrls: ['./sliding-panel.component.scss'],
  standalone: true
})
export class SlidingPanelComponent {

  @Input() visible : boolean = false;
  @Output() visibleChange = new EventEmitter<boolean>();
  @Input() title: string = ""

  constructor() {}

  showPanel(event: any){
    event.stopPropagation();
  }

  closePanel(){
    this.visible = false
    this.visibleChange.emit(this.visible);
  }

  @HostListener('document:keydown.escape', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent) {
    if (this.visible) {
      this.closePanel();
    }
  }
}
