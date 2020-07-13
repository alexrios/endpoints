# Endpoints

The tiniest http endpoints simulator

## Status
This project is under development.

### Usage
On terminal
```shell script
$ endpoints 
```
You should see 
```
INFO[0000] [POST] / -> 200 with body -> customBody.json 
INFO[0000] Listen at :8080                              
```

#### Defaults
* address - ":8080"
* method - "GET"
* latency - 0ms
* status code - 200

#### Configuration file
####### endpoints.json

All features example:
```json
{
  "address": ":8080",
  "responses": [
    {
      "path": "/",
      "status": 201,
      "latency": "400ms",
      "method": "POST",
      "json_body": "customBody.json"
    }
  ]
}
```

note: json_body is the file location of body file.