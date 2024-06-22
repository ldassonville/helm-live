import { Component, OnInit } from '@angular/core';
import { SchemaService } from '../schema.service';
import { JsonPipe, KeyValuePipe, NgFor } from '@angular/common';
import { PropertyViewerComponent } from '../property-viewer/property-viewer.component';
import { NgModel } from '@angular/forms';
import { SchemaResolver } from '../schema-resolver';

@Component({
  selector: 'app-schema-viewer',
  templateUrl: './schema-viewer.component.html',
  styleUrl: './schema-viewer.component.scss'
})
export class SchemaViewerComponent implements OnInit {

  protected schema: any

  protected schemasUrls: string[] = []
  protected schemaUrl: string = ""



  constructor(private schemaResolver: SchemaResolver){

      this.schemasUrls = [
        "./assets/schemas-catalog/helm/helm-api.schema.json",
"./assets/schemas-catalog/externaldns.k8s.io/dnsendpoint_v1alpha1.json",
"./assets/schemas-catalog/cdn.aap.adeo.cloud/xrccloudflareservice_v1alpha1.json",
"./assets/schemas-catalog/dsp.aap.adeo.cloud/xrctopic_v1alpha1.json",
"./assets/schemas-catalog/dsp.aap.adeo.cloud/xrcreference_v1alpha1.json",
"./assets/schemas-catalog/dsp.aap.adeo.cloud/xrcapikey_v1alpha1.json",
"./assets/schemas-catalog/tf.upbound.io/workspace_v1beta1.json",
"./assets/schemas-catalog/generators.external-secrets.io/password_v1alpha1.json",
"./assets/schemas-catalog/aap.adeo.cloud/xrcpgdatabase_v1alpha1.json",
"./assets/schemas-catalog/aap.adeo.cloud/xrcgkecluster_v1alpha1.json",
"./assets/schemas-catalog/aap.adeo.cloud/xrcpgservice_v1alpha1.json",
"./assets/schemas-catalog/aap.adeo.cloud/xrcdbapipg_v1alpha1.json",
"./assets/schemas-catalog/configuration.konghq.com/kongconsumer_v1.json",
"./assets/schemas-catalog/configuration.konghq.com/kongplugin_v1.json",
"./assets/schemas-catalog/cert-manager.io/certificate_v1.json",
"./assets/schemas-catalog/external-secrets.io/externalsecret_v1beta1.json",
"./assets/schemas-catalog/external-secrets.io/pushsecret_v1alpha1.json",
"./assets/schemas-catalog/external-secrets.io/externalsecret_v1alpha1.json",
"./assets/schemas-catalog/datadoghq.com/watermarkpodautoscaler_v1alpha1.json",
"./assets/schemas-catalog/argoproj.io/application_v1alpha1.json",


      ]
      this.schemaUrl = this.schemasUrls[0]

  }


  ngOnInit(): void {
    
    this.loadSchema()
  }

  selectSchema($event: any){
    this.schemaUrl = $event.target.value
    this.loadSchema()
    
  }


  loadSchema(){
    this.schema = undefined
    console.log("loadSchema")
    this.schemaResolver
      .init(this.schemaUrl)
      .subscribe(schema => {
        this.schema = schema;
      }
    )
  }
}
