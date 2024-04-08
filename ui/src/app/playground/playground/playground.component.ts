import { Component} from '@angular/core';
import {HelmRender} from "../render";
import {RenderService} from "../render.service";
import { JsonPipe} from "@angular/common";
import {MonacoEditorModule, NGX_MONACO_EDITOR_CONFIG} from "ngx-monaco-editor-v2";
import {FormsModule} from "@angular/forms";


@Component({
  selector: 'app-playground',
  standalone: true,
  imports: [
    JsonPipe,
    MonacoEditorModule,
    FormsModule
  ],
  providers: [{ provide: NGX_MONACO_EDITOR_CONFIG, useValue: NGX_MONACO_EDITOR_CONFIG }],
  templateUrl: './playground.component.html',
  styleUrl: './playground.component.scss',
})
export class PlaygroundComponent{

  public render: HelmRender = {}

  protected editorOptions = {theme: 'vs', language: 'yaml', fontSize: 13.5, automaticLayout: true};
  protected content = ""

  protected selection = {
    category: "values",
    identity: "mergedValues",
  }

  constructor( private renderService: RenderService) {
    this.renderService.
      live().
      subscribe(render => {
        this.render = render;
        this.selectItem(this.selection.category, this.selection.identity)
    });
  }

  onSelectManifest(group: string, kind: string, version: string, name: string) {

    this.selection = {
      category: "manifests",
      identity: `${group}-${kind}-${version}-${name}`
    }
    this.selectItem(this.selection.category, this.selection.identity)
  }

  onSelectValues(identifier: string) {
    this.selection = {
      category: "values",
      identity: identifier
    }
    this.selectItem(this.selection.category, this.selection.identity)
  }

  selectManifestByIdentifier(identifier: string) {
    console.log("selectManifestByIdentifier", identifier)
    if(this.render.manifests){
      const foundManifest = this.render.manifests?.find(manifest => {
        return <string>`${manifest.group}-${manifest.kind}-${manifest.version}-${manifest.name}` === identifier
      })
      this.content = foundManifest?.content || ""
    }else {
      this.content = ""
    }
  }



  selectItem(category: string, identifier: string) {

    switch (category) {
      case "manifests":
        this.selectManifestByIdentifier(identifier)
        break;
      case "values":
        if(identifier === "mergedValues") {
          this.content = this.render?.mergedValues?.data || ""
        }
        break;
      default:
        this.content = ""
    }
  }
}
