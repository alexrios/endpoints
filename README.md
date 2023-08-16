The tiniest HTTP endpoints simulator

[![Go Report Card](https://goreportcard.com/badge/github.com/alexrios/endpoints)](https://goreportcard.com/report/github.com/alexrios/endpoints)
[![Download shield](https://img.shields.io/github/downloads/alexrios/endpoints/total)](https://img.shields.io/github/downloads/alexrios/endpoints/total)


# Endpoints

### Usage
Download your binary [here](https://github.com/alexrios/endpoints/releases/latest)


On terminal
```shell script
$ endpoints 
```
You should see 
```
INFO[0000] [GET] / -> 200 with body -> customBody.json 
INFO[0000] Listen at :8080                              
```

#### Defaults
* address - ":8080"
* method - "GET"
* latency - 0ms
* status code - 200

#### Configuration file
##### endpoints.json

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

Note: json_body is the file location of the body file.

##### Body interpolation with path variables
Now you wanna interpolate an identifier on the response body. How to do it?

Let's add another response on `responses`.
```json
{
  "address": ":8080",
  "responses": [
    {
      "path": "/",
      "status": 200,
      "latency": "400ms",
      "method": "POST",
      "json_body": "customBody.json"
    },
    {
      "path": "/{id}/sales",
      "status": 201,
      "latency": "400ms",
      "method": "GET",
      "json_body": "interpolated.json"
    }
  ]
}
```
And now, we'll use templating notation to use this response body as a template.
##### interpolated.json
```
{
  "id": {{ .id}}
}
```

## Status
This project is under development.
