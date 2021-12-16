package main

import (
	_ "os"
	"fmt"
	"strconv"
	"net/http"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	_ "github.com/amleshkashyap/respTest/generator"
)


func readAndReturnJson(endpoint string) (string, []byte, bool, string) {
	endpointData, err := ioutil.ReadFile("examples/endpoints.json")
        if err != nil {
                return "", nil, false, "Can't execute without input file"
        }

        endpoints := make(map[string][]int)
        json.Unmarshal(endpointData, &endpoints)


	var statusCodes []int
	for key, val := range endpoints {
		if key == endpoint {
			statusCodes = val
		}
	}

	if len(statusCodes) == 0 {
		return "", nil, false, "Endpoint data not found."
	}

	randomIndex := rand.Intn(len(statusCodes))
	code := fmt.Sprintf("%d", statusCodes[randomIndex])

	fmt.Printf("Server endpoint: %s, returning code: %s\n", endpoint, code)

	responseJson, err := ioutil.ReadFile("examples/data.json")
	if err != nil {
		return "", nil, false, fmt.Sprintf("Error occurred in reading file: %s\n", err.Error())
	}

	var someData json.RawMessage
	json.Unmarshal(responseJson, someData)

	var data map[string]json.RawMessage
	json.Unmarshal(responseJson, &data)

	if data[endpoint] == nil {
		return "", nil, false, "Data for this endpoint isn't present, returning."
	}

	var thisData map[string]json.RawMessage
	json.Unmarshal(data[endpoint], &thisData)

	if thisData[code] == nil {
		fmt.Printf("Data for this error code isn't present, finding the first error code with data.\n")
		codeTemp := ""
		for k, _ := range thisData {
			codeTemp = k
			break
		}
		if codeTemp == code {
			return "", nil, false, "No error code specific data present for this endpoint, returning."
		}
		code = codeTemp
	}

	var actualData map[string]json.RawMessage
	json.Unmarshal(thisData[code], &actualData)

	var keyArr []string
	for k, _ := range actualData {
		keyArr = append(keyArr, k)
	}

	if len(keyArr) == 0 {
		return "", nil, false, "No actual data present, returning."
	}

	randomIndex = rand.Intn(len(keyArr))
	fmt.Printf("Returning value for endpoint: %s, code: %s, subcode: %s\n", endpoint, code, keyArr[randomIndex])

	return code, actualData[keyArr[randomIndex]], true, ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("Full path is %s\n", path)
	// setup the content-type before writing any status codes
	w.Header().Set("Content-Type", "application/json")
	code, data, success, message := readAndReturnJson(path)
	if success == false {
	        fmt.Printf("Error Occurred: %s\n", message)
	        w.WriteHeader(599)
	} else {
	        codeVal, err := strconv.Atoi(code)
	        if err != nil {
	                fmt.Println("Something went wrong")
	                w.WriteHeader(599)
	        } else {
	                w.WriteHeader(codeVal)
	                w.Write(data)
	        }
	}
}

func httpServer() {
	http.HandleFunc("/", handler)
        fmt.Println("Starting HTTP Server on Port: ", 3000)
	err := http.ListenAndServe(":3000", nil)
        if err != nil {
                panic(err.Error())
        }
}

func main() {
	httpServer()
}
