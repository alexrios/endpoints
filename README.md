# Endpoints

The tiniest http endpoints simulator

## Status
This project is under development. No use this in production :)


#### 

##### endpoints.json
```json
{
  "address": ":8080",
  "responses": [
    {
      "path": "/backend/test/{name}",
      "status": 201,
      "latency": "400ms",
      "json_body": "back.json"
    }
  ]
}
```

##### JSON body
eg.: back.json
```json
{
  "nome": "{{ .name}}"
}
```