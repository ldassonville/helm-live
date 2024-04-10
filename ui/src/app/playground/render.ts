export interface HelmRender {

  info: RenderInfo,

  chartPath? :string;
  chartName?: string;

  rawManifest?: string;

  sources?: HelmSourceFile[];
  errors?: RenderError[];

  values?: Values
  valuesFiles? : Values[];

  mergedValues?: Values;

}

export interface HelmSourceFile {
  source: string;
  manifests?: Manifest[];
}

export interface KubeConformValidation  {
  status: string
  errMsg?: string
  validationErrors?: ValidationError[]
}

export interface ValidationError {
  path: string;
  message: string;
}

export interface  RenderInfo {
  status?: string
  executionTime?: Date
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
  source: string;

  // Kubernetes validation
  kubeConformValidation?: KubeConformValidation

  // Yaml validation
  isYamlValid: boolean;
  yamlError: string
}


export interface  ValidationError {
  path: string;
  msg: string;
}

export interface RenderError {
  message: string;
}
