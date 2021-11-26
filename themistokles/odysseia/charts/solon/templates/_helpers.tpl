{{/*
Expand the name of the chart.
*/}}
{{- define "solon.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "solon.fullname" -}}
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
{{- define "solon.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "solon.labels" -}}
helm.sh/chart: {{ include "solon.chart" . }}
{{ include "solon.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "solon.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "solon.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the role to use
*/}}
{{- define "solon.roleName" -}}
{{- if .Values.role.create -}}
{{- default (include "solon.fullname" .) .Values.role.name }}
{{- else -}}
    {{ default "default" .Values.role.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the roleBinding to use
*/}}
{{- define "solon.bindingName" -}}
{{- if .Values.role.create -}}
{{- default (include "solon.fullname" .) .Values.role.name }}
{{- else -}}
    {{- default (include "solon.roleName" .) "binding" }}
{{- end -}}
{{- end -}}

{{- define "solon.uname" -}}
{{- if empty .Values.services.solon.name -}}
{{ .Values.images.solon }}
{{- else -}}
{{ .Values.services.solon.name }}
{{- end -}}
{{- end -}}

{{/*
Allow the release namespace to be overridden for multi-namespace deployments in combined charts
*/}}
{{- define "solon.namespace" -}}
  {{- if .Values.namespaceOverride -}}
    {{- .Values.namespaceOverride -}}
  {{- else -}}
    {{- .Release.Namespace -}}
  {{- end -}}
{{- end -}}
