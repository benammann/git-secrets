# Created by git-secrets
# Context: {{ .UsedContext }}
# Origin File {{ .UsedFile.FileIn }}
# Destination File {{ .UsedFile.FileOut }}

apiVersion: v1
data:
  password: {{ Base64Encode .Secrets.databasePassword  }}
kind: Secret
metadata:
  name: database-{{ .UsedContext }}
type: Opaque