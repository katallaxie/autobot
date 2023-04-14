{{- define "autobot.namespace" -}}
{{- printf "%s" .Values.namespace }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "autobot.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "autobot.fullname" -}}
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
Expand the name of the chart.
*/}}
{{- define "autobot.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "autobot.labels" -}}
helm.sh/chart: {{ include "autobot.chart" . }}
{{ include "autobot.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "autobot.selectorLabels" -}}
app.kubernetes.io/name: {{ include "autobot.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "autobot.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "autobot.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account secret token to use
*/}}
{{- define "autobot.serviceAccountTokenName" -}}
{{- if .Values.serviceAccount.create }}
{{- printf "%s-%s" (include "autobot.fullname" .) .Values.serviceAccount.tokenName }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "autobot.env" -}}
{{- with .Values.env }}
    {{- range $ev_key, $ev_value :=  . }}
        {{- if (typeIs "string" $ev_value) }}
- name: {{ $ev_key }}
  value: {{ $ev_value | quote }}
        {{- else }}
- name: {{ $ev_key }}
{{- toYaml $ev_value | nindent 2 }}
        {{- end }}
    {{- end }}
{{- end }}
{{- end -}}

{{/*
Create a default fully qualified app rootname.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
Values.global.rootname will be assigned in laas-deployment ApplicationSet
*/}}
{{- define "autobot.rootname" -}}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if .Values.global }}
{{- if .Values.global.rootname }}
{{- printf "%s.%s" $name .Values.global.rootname | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- else }}
{{- $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}