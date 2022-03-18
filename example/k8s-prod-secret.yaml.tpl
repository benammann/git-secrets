apiVersion: v1
data:
  password: {{ base64 .Secrets.databasePassword }}
kind: Secret
metadata:
  name: database
type: Opaque