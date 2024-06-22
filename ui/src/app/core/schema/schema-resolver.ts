import { Observable, catchError, map, of, throwError } from "rxjs";
import { SchemaService } from "./schema.service";
import { Injectable } from "@angular/core";
import urlJoin from 'url-join';


@Injectable({
  providedIn: 'root'
})
export class SchemaResolver {

  rootSchema: string = ""
  schemaWorkdir: string = ""

  schemaCache = new Map<string, string>();

  /**
   * @constructor
   */
  constructor(private schemaClient: SchemaService) {
  }

  public init(schemaPath: string): Observable<any> {
    this.rootSchema = schemaPath

    return this.schemaClient
      .getSchema(this.rootSchema)
      .pipe(
        map(schema => this.registerToCache(this.rootSchema, schema)),
        map(schema => {
          this.schemaWorkdir = new SchemaPath(schemaPath).path
          return schema
        })
      )
  }

  private findByPath(path: string, schema: any): any {

    let currentSchema = schema
    if(path){
      let parts = path.split("/")
      for(let i=0; i < parts.length; i++){

        const eltName = parts[i]
        if(! eltName) {
          continue
        }
        currentSchema = currentSchema[eltName]

        if(!currentSchema){
          break
        }
      }
    }
    return currentSchema
  }

  public resolveTarget(path: string, schema: any): {schema: any, target: any} {

    console.log("resolveTarget")
    let target: any = this.findByPath(path, schema)

    return {
      schema: schema, target: target
    }
  }

  public registerToCache(schemaURL: string, schema: any): any {
    if (schema) {
      this.schemaCache.set(schemaURL, schema)
    }
    return schema
  }

  public getAbsoluteURL(currentFile: string, refStr: string): string {
    let ref = new SchemaRef(refStr, currentFile)

    if(!ref.schemaPath ){
      return currentFile
    }

    const currentFileSchemaPath = new SchemaPath(currentFile)

    console.log("getAbsoluteURL")
    let absoluteURL = urlJoin(
      //this.schemaWorkdir,
      currentFileSchemaPath.path,
      ref.schemaPath.filepath ? ref.schemaPath.filepath : currentFileSchemaPath.filename
    )

    return absoluteURL
  }

  public resolveRef(currentFile: string, refStr: string): Observable<{schema: any, target: any}> {

    const absoluteURL = this.getAbsoluteURL(currentFile, refStr)
    const elementPath = new SchemaRef(refStr, currentFile).elementPath

    if (this.schemaCache.has(absoluteURL)) {
      console.log("Using caches schema");

      const schema = this.schemaCache.get(absoluteURL)

      return of( {
        schema: this.schemaCache.get(absoluteURL),
        target: this.findByPath(elementPath, schema)
      })
    }

    return this.schemaClient
      .getSchema(absoluteURL)
      .pipe(
        catchError(err => {
          if (err.error.statusCode == 404) {
            return of(undefined)
          } else {
            return throwError(() => err)
          }
        }),
        map(schema => this.registerToCache(absoluteURL, schema)),
        map(schema => this.resolveTarget(elementPath, schema)),
      )
  }
}


export class SchemaRef {

  private _strRef: string
  private _srcFile: string
  private parsedRef: { schemaPath: SchemaPath, elementPath: string }

  constructor(strRef: string, srcFile: string) {

    this._strRef = strRef
    this._srcFile = srcFile

    // Analyse the given ref
    this.parsedRef = this.parseStrRef(strRef)
  }


  private parseStrRef(strRef: string) {

    let parts = strRef.split("#")
    return {
      schemaPath: new SchemaPath(parts[0]),
      elementPath: parts[1]
    }
  }

  get ref(): string {
    return this._strRef
  }

  get srcFile(): string {
    return this._srcFile
  }

  get schemaPath(): SchemaPath {
    return this.parsedRef.schemaPath
  }
  get elementPath(): string {
    return this.parsedRef.elementPath
  }

}


export class SchemaPath {


  private _filepath: string = ""
  private parsedPath: { path: string, filename: string }

  /**
   * @param filepath
   * Ex: schemas/product/components/api.schema.json
   */
  constructor(filepath: string) {
    this._filepath = filepath
    this.parsedPath = this.parsePath(filepath)
  }

  /**
   * Path a file path
   * @param path to analyse
   * @returns Passered element
   *
   * Ex schemas/sample/path/demo.schema.json
   * {
   *    "path" : "schemas/sample/path/"
   *    "file" : "demo.schema.json"
   * }
   */
  private parsePath(path: string) {
    let i = path.lastIndexOf("/");

    return {
      path: path.substring(0, i),
      filename: path.substring(i)
    }
  }

  /*
    Return the filepath
    ex: schemas/product/components/api.schema.json
  */
  public get filepath(): string {
    return this._filepath
  }

  /*
    Return the filepath
    ex: schemas/product/components
  */
  public get path(): string {
    return this.parsedPath.path
  }

  /*
    Return the filename
    ex: api.schema.json
  */
  public get filename(): string {
    return this.parsedPath.filename
  }
}
