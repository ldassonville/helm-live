import {Component} from '@angular/core';
import {HelmRender, KubeConformValidation, Manifest, ValidationError} from "../render";
import {RenderService} from "../render.service";
import {JsonPipe} from "@angular/common";
import {MonacoEditorModule, NGX_MONACO_EDITOR_CONFIG} from "ngx-monaco-editor-v2";
import {FormsModule} from "@angular/forms";
import {SlidingPanelComponent} from "../../core/sliding-panel/sliding-panel.component";
import {SchemaViewerComponent} from "../schema-viewer/schema-viewer.component";

@Component({
  selector: 'app-playground',
  standalone: true,
  imports: [
    JsonPipe,
    MonacoEditorModule,
    FormsModule,
    SlidingPanelComponent,
    SchemaViewerComponent,
  ],
  providers: [{provide: NGX_MONACO_EDITOR_CONFIG, useValue: NGX_MONACO_EDITOR_CONFIG}],
  templateUrl: './playground.component.html',
  styleUrl: './playground.component.scss',
})
export class PlaygroundComponent {

  public render: HelmRender = {info: {}}

  protected kubeConformValidation : KubeConformValidation|null = null
  protected editorOptions = {theme: 'vs', language: 'yaml', fontSize: 12.5, automaticLayout: true, readOnly: true};
  protected content = ""
  protected editor: any
  protected showResourceDetailPanel = false

  protected selectedGVK: { group: string, kind: string, version: string }| undefined
  protected sourcesSelection = new Map<string, boolean>();

  protected selection = {
    category: "values",
    identity: "mergedValues",
  }

  constructor(private renderService: RenderService) {

    this.renderService.render().subscribe(render => {
      this.render = render;
    })

    this.renderService.
      live().
      subscribe(render => {
        this.render = render;
        this.selectItem(this.selection.category, this.selection.identity)
    });
  }

  onEditorInit(editor: any) {
    this.editor = editor
  }

  onSelectRawManifest() {
    this.selection = {
      category: "rawManifest",
      identity: "rawManifest"
    }
    this.selectItem(this.selection.category, this.selection.identity)
  }

  onSelectSource(source: string) {
    let selected = this.isSelectedSource(source)
    this.sourcesSelection.set(source, !selected)
  }

  isSelectedSource(source: string): boolean {
    let selected = true
    if (this.sourcesSelection.has(source)) {
      selected = <boolean>this.sourcesSelection.get(source)
    }
    return selected
  }

  onSelectManifest(group: string, kind: string, version: string, name: string) {

    this.selection = {
      category: "manifests",
      identity: `${group}-${kind}-${version}-${name}`
    }
    if (this.editor) {
      this.editor.setScrollPosition({scrollTop: 0});
      //editor.revealLine(15);
    }

    this.selectItem(this.selection.category, this.selection.identity)
  }

  isSelectedManifest(group: string, kind: string, version: string, name: string): boolean {

    if (this.selection.category !== "manifests") {
      return false
    }

    return this.selection.identity == `${group}-${kind}-${version}-${name}`
  }

  onSelectValues(identifier: string) {
    this.selection = {
      category: "values",
      identity: identifier
    }
    this.selectItem(this.selection.category, this.selection.identity)
  }

  selectManifestByIdentifier(identifier: string) {

    if (this.render.sources) {
      let foundManifest: any = undefined
      for (let source of this.render.sources) {
        foundManifest = source.manifests?.find(manifest => {
          return <string>`${manifest.group}-${manifest.kind}-${manifest.version}-${manifest.name}` === identifier
        })
        if (foundManifest) {
          break
        }
      }
      this.content = foundManifest?.content || ""
      this.kubeConformValidation = foundManifest?.kubeConformValidation || null
    } else {
      this.content = ""
    }
  }



  selectItem(category: string, identifier: string) {
    this.kubeConformValidation = null
    switch (category) {
      case "rawManifest":
        this.content = this.render.rawManifest || ""
        break;
      case "manifests":
        this.selectManifestByIdentifier(identifier)
        break;
      case "values":
        if (identifier === "mergedValues") {
          this.content = this.render?.mergedValues?.data || ""
        }
        break;
      default:
        this.content = ""
    }
  }

  showKindDetail(manifest: Manifest) {
    this.showResourceDetailPanel = true
    this.selectedGVK = {
      group: manifest.group,
      kind: manifest.kind,
      version: manifest.version,
    }
  }
}
