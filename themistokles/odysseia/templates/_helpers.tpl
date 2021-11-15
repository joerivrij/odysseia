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

{{/*
Create the name of the service account to use
*/}}
{{- define "odysseia.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "odysseia.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "hippokrates.uname" -}}
{{- if empty .Values.services.hippokrates.name -}}
{{ .Values.images.hippokrates }}
{{- else -}}
{{ .Values.services.hippokrates.name }}
{{- end -}}
{{- end -}}

{{- define "alexandros.uname" -}}
{{- if empty .Values.services.alexandros.name -}}
{{ .Values.images.alexandros }}
{{- else -}}
{{ .Values.services.alexandros.name }}
{{- end -}}
{{- end -}}

{{- define "dionysos.uname" -}}
{{- if empty .Values.services.dionysos.name -}}
{{ .Values.images.dionysos }}
{{- else -}}
{{ .Values.services.dionysos.name }}
{{- end -}}
{{- end -}}

{{- define "herodotos.uname" -}}
{{- if empty .Values.services.herodotos.name -}}
{{ .Values.images.herodotos }}
{{- else -}}
{{ .Values.services.herodotos.name }}
{{- end -}}
{{- end -}}

{{- define "sokrates.uname" -}}
{{- if empty .Values.services.sokrates.name -}}
{{ .Values.images.sokrates }}
{{- else -}}
{{ .Values.services.sokrates.name }}
{{- end -}}
{{- end -}}

{{- define "solon.uname" -}}
{{- if empty .Values.services.solon.name -}}
{{ .Values.images.solon }}
{{- else -}}
{{ .Values.services.solon.name }}
{{- end -}}
{{- end -}}

{{- define "ptolemaios.uname" -}}
{{- if empty .Values.services.ptolemaios.name -}}
{{ .Values.images.ptolemaios }}
{{- else -}}
{{ .Values.services.ptolemaios.name }}
{{- end -}}
{{- end -}}