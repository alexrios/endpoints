package main

//Configuration file name.
const DefaultConfigurationFileName = "endpoints.json"

//Configuration file content to create an example on very first run.
const DefaultConfigurationFileContent = `{
  "responses": [
    {
      "path": "/",
      "json_body": "customBody.json"
    }
  ]
}

`

const CustomBodyExampleFileName = "customBody.json"

//Customized json body content to create an example on very first run.
const CustomBodyExampleFileContent = `{
  "name": "Alex"
}`
