{{/*
Expand the name of the chart.
*/}}
{{- define "odysseia.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "odysseia.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "odysseia.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "odysseia.labels" -}}
helm.sh/chart: {{ include "odysseia.chart" . }}
{{ include "odysseia.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "odysseia.selectorLabels" -}}
app.kubernetes.io/name: {{ include "odysseia.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "odysseiaapi.uname" -}}
{{- if empty .Values.services.odysseiaapi.name -}}
{{ .Values.images.odysseiaapi }}
{{- else -}}
{{ .Chart.Name }}
{{- end -}}
{{- end -}}

{{- define "sidecar.uname" -}}
{{- if empty .Values.services.sidecar.name -}}
{{ .Values.images.sidecar }}
{{- else -}}
{{ .Values.services.sidecar.name }}
{{- end -}}
{{- end -}}

{{- define "init.uname" -}}
{{- if empty .Values.services.init.name -}}
{{ .Values.images.init }}
{{- else -}}
{{ .Values.services.init.name }}
{{- end -}}
{{- end -}}

{{- define "seeder.uname" -}}
{{- if empty .Values.services.seeder.name -}}
{{ .Values.images.seeder }}
{{- else -}}
{{ .Values.services.seeder.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "serviceAccountName" -}}
  {{- if .Values.serviceAccountOverride -}}
    {{- .Values.serviceAccountOverride -}}
  {{- else -}}
    {{- .Chart.Name -}}
  {{- end -}}
{{- end -}}

{{/*
Create the name of the role to use
*/}}
{{- define "roleName" -}}
  {{- if .Values.roleNameOverride -}}
    {{- .Values.roleNameOverride -}}
  {{- else -}}
    {{- .Chart.Name -}}
  {{- end -}}
{{- end -}}

{{/*
Create the name of the roleBinding to use
*/}}
{{- define "bindingName" -}}
{{- default "dionysos-binding" }}
{{- end -}}

{{/*
Allow the release namespace to be overridden for multi-namespace deployments in combined charts
*/}}
{{- define "namespace" -}}
  {{- if .Values.namespaceOverride -}}
    {{- .Values.namespaceOverride -}}
  {{- else -}}
    {{- .Release.Namespace -}}
  {{- end -}}
{{- end -}}