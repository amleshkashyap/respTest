package main

import (
	"os"
	"fmt"
	"strconv"
	"strings"
	"regexp"
	"net/http"
	"math/rand"
	"io/ioutil"
	"encoding/json"
)

var MODE string

func readAndReturnMode1(endpoint string) (string, []byte, bool, string) {
	prefix := ""
	if len(os.Args) > 2 {
		prefix = "samples"
	} else {
		prefix = "examples"
	}

	endpointData, err := ioutil.ReadFile(fmt.Sprintf("%s/endpoints.json", prefix))
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

	responseJson, err := ioutil.ReadFile(fmt.Sprintf("%s/data.json", prefix))
	if err != nil {
		return "", nil, false, fmt.Sprintf("Error occurred in reading file: %s\n", err.Error())
	}

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

func readAndReturnMode2(endpoint, param string) (int, []byte, map[string]string, bool, string) {
        prefix := ""
        if len(os.Args) > 2 {
                prefix = "samples"
        } else {
                prefix = "examples"
        }

        endpointData, err := ioutil.ReadFile(fmt.Sprintf("%s/new_data.json", prefix))
        if err != nil {
                return 599, nil, nil, false, "Can't execute without input file"
        }

        var endpoints map[string]json.RawMessage
        json.Unmarshal(endpointData, &endpoints)

	if endpoints[endpoint] == nil {
		found := false
		for key, _ := range endpoints {
			if strings.Contains(key, "{id}") {
				newEndpointStr := strings.Replace(key, "{id}", ".*", 1)
				matched, _ := regexp.MatchString(newEndpointStr, endpoint)
				if matched == true {
					found = true
					endpoint = key
					break
				}
			}
		}
		if found == false {
			return 599, nil, nil, false, "Data for this endpoint isn't present, returning."
		}
	}

	var paramData map[string]json.RawMessage
	json.Unmarshal(endpoints[endpoint], &paramData)

	if paramData[param] == nil {
		return 599, nil, nil, false, "Data for this unique parameter isn't present, returning"
	}

	var data map[string]json.RawMessage
	json.Unmarshal(paramData[param], &data)

	if (data["code"] != nil && data["response"] != nil) {
		var code int
		headers := make(map[string]string)
		json.Unmarshal(data["code"], &code)
		json.Unmarshal(data["headers"], &headers)
		return code, data["response"], headers, true, "Success"
	}
	return 599, nil, nil, false, "Couldn't find endpoint data"
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("Full path is %s\n", path)
	paramId := r.URL.Query().Get(os.Getenv("UNIQUE_PARAM_ID"))

	if MODE == "mode-2" {
		code, data, headers, success, message := readAndReturnMode2(path, paramId)
		if success == false {
			fmt.Printf("Error Occurred: %s\n", message)
			w.WriteHeader(code)
		} else {
			for key, val := range headers {
				w.Header().Set(key, val)
			}
			w.WriteHeader(code)
			w.Write(data)
		}
	} else {
		code, data, success, message := readAndReturnMode1(path)
		w.Header().Set("Content-Type", "application/json")
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
}

func httpServer() {
	http.HandleFunc("/", handler)
	if len(os.Args) < 2 {
		fmt.Println("Can't run the server without the mode")
		os.Exit(1)
	}
	MODE = os.Args[1]
        fmt.Println("Starting HTTP Server on Port: ", 3000)
	err := http.ListenAndServe(":3000", nil)
        if err != nil {
                panic(err.Error())
        }
}

func main() {
	httpServer()
}
