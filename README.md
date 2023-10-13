# portainerexport
 Local tool to connect to a several Portainer.io API and collect information about versions that is running there.

 config.json.example - contains configuration example
 keys ar written in config file as there is in local runtime you will need to store some keys. Possible improvement to store them encrypted

 if you run from code - just run `go run main.go` - it will generate an output.xlsx - that will contain result as containers that was found in each enviornment and will write version that was found.

 Output format can be customized with a flag:
 
 `--format string   Format in which to render output. Supported formats: excel, markdown (default "excel")`

 Possible improvement - if there is no container found then also add line in result as - not found result. Now it does not write such a line.
