package processing

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	model "github.com/zviedris/portainerexport/model"
)

func callPortainer(env model.Enviornment, client *http.Client, config *model.Config, useStack int16, stack model.Stack, results map[string][]model.EnvVersion) map[string][]model.EnvVersion {
	req, err := http.NewRequest("GET", env.Url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return results
	}

	// Set the API key in the request header
	req.Header.Add("X-API-Key", env.ApiKey)

	//if filter by a stack then
	// Add query parameter
	if useStack == 1 {
		q := req.URL.Query()

		filterValue := "{\"label\":[\"com.docker.stack.namespace=" + stack.Name + "\"]}"

		q.Add("filters", filterValue)
		req.URL.RawQuery = q.Encode()
	}

	// Send the GET request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return results
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("Request %s failed with status code: %d\n", req.URL, response.StatusCode)
		return results
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return results
	}

	var portList []model.PortObject

	err = json.Unmarshal(body, &portList)
	if err != nil {
		fmt.Println("Error:", err)
		return results
	}

	for _, cont := range portList {
		//key := cont.Spec.Name
		excludeVal := false
		for _, excl := range config.Exclude {
			if strings.Contains(cont.Spec.Labels.Image, excl.Name) {
				excludeVal = true
			}
		}
		if !excludeVal {
			images := strings.Split(cont.Spec.Labels.Image, ":")
			path := strings.Split(images[0], "/")
			key := path[len(path)-1]
			var value model.EnvVersion
			value.Environment = env.Name
			value.Docker = images[1]
			value.Stack = cont.Spec.Labels.Namespace
			value.DockerPath = images[0]
			existingVal := results[key]
			existingVal = append(existingVal, value)
			results[key] = existingVal
		}

	}
	return results
}

func ProcessPortainer(config *model.Config) map[string][]model.EnvVersion {

	results := map[string][]model.EnvVersion{}

	for _, env := range config.Environments {

		//ignore that Portainer need verified certificate
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// Create a new HTTP client
		client := &http.Client{Transport: tr}

		if config.UseStacks == 1 {
			//for each stack
			for _, stack := range config.Stacks {
				// Create a new GET request with query parameter
				results = callPortainer(env, client, config, 1, stack, results)
			}
		} else {
			results = callPortainer(env, client, config, 0, model.Stack{}, results)
		}

	}

	return results
}
