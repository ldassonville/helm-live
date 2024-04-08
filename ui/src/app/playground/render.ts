export interface HelmRender {

  manifests?: Manifest[];
  errors?: RenderError[];

  values?: Values
  valuesFiles? : Values[];

  mergedValues?: Values;
}


export interface Values {
  data: string;
}

export interface GVK {
  group: string;
  version: string;
  kind: string;
}

export interface Manifest extends GVK{
  name: string;
  namespace: string;
  content: string;
}

export interface RenderError {
  message: string;
}
