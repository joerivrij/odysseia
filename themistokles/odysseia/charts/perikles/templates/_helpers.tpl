{{/*
Expand the name of the chart.
*/}}
{{- define "perikles.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "perikles.fullname" -}}
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
{{- define "perikles.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "perikles.labels" -}}
helm.sh/chart: {{ include "perikles.chart" . }}
{{ include "perikles.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "perikles.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "perikles.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the role to use
*/}}
{{- define "perikles.roleName" -}}
{{- if .Values.role.create -}}
{{- default (include "perikles.fullname" .) .Values.role.name }}
{{- else -}}
    {{ default "default" .Values.role.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the roleBinding to use
*/}}
{{- define "perikles.bindingName" -}}
{{- if .Values.role.create -}}
{{- default (include "perikles.fullname" .) .Values.role.name }}
{{- else -}}
    {{- default (include "perikles.roleName" .) "api-access-binding" }}
{{- end -}}
{{- end -}}

{{- define "perikles.uname" -}}
{{- if empty .Values.services.perikles.name -}}
{{ .Values.images.perikles }}
{{- else -}}
{{ .Values.services.perikles.name }}
{{- end -}}
{{- end -}}

{{/*
Allow the release namespace to be overridden for multi-namespace deployments in combined charts
*/}}
{{- define "perikles.namespace" -}}
  {{- if .Values.namespaceOverride -}}
    {{- .Values.namespaceOverride -}}
  {{- else -}}
    {{- .Release.Namespace -}}
  {{- end -}}
{{- end -}}


{{- define "drakon.uname" -}}
{{- if empty .Values.services.drakon.name -}}
{{ .Values.images.drakon }}
{{- else -}}
{{ .Values.services.drakon.name }}
{{- end -}}
{{- end -}}
