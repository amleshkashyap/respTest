package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"regexp"
	"net/http"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"github.com/amleshkashyap/respTest/compute"
)

var MODE string
var MATCHER string
var DATAPATH string

func readAndReturnMode1(endpoint string) (string, []byte, bool, string) {
	endpointData, err := ioutil.ReadFile(fmt.Sprintf("%s/endpoints.json", DATAPATH))
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

	responseJson, err := ioutil.ReadFile(fmt.Sprintf("%s/data.json", DATAPATH))
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
        endpointData, err := ioutil.ReadFile(fmt.Sprintf("%s/new_data.json", DATAPATH))
        if err != nil {
                return 599, nil, nil, false, "Can't execute without input file"
        }

        var endpoints map[string]json.RawMessage
        json.Unmarshal(endpointData, &endpoints)

	if endpoints[endpoint] == nil {
		found, max, ePoint := false, 0, ""
		if MATCHER == "Custom" {
			for key, _ := range endpoints {
				if strings.Contains(key, "{id}") == false { continue }
				endpointStr := strings.Replace(key, "{id}", "", 1)
				length, match := compute.UrlPatternMatch(endpointStr, endpoint)
				if match == false { continue }
				found = true
				if length > max {
					max = length
					ePoint = key
				}
			}
		} else {
			for key, _ := range endpoints {
				if strings.Contains(key, "{id}") {
					newEndpointStr := strings.Replace(key, "{id}", ".*", 1)
					matched, _ := regexp.MatchString(newEndpointStr, endpoint)
					if matched == true {
						found = true
						endpointStr := strings.Replace(key, "{id}", "", 1)
						result := compute.Longest(endpointStr, endpoint)
						if len(result) > max {
							max = len(result)
							ePoint = key
						}
					}
				}
			}
		}
		if found == false {
			return 599, nil, nil, false, "Data for this endpoint isn't present, returning."
		}
		endpoint = ePoint
	}

	// fmt.Printf("    Serving from selected endpoint: %s, for param: %s\n\n", endpoint, param)

	var paramData map[string]json.RawMessage
	json.Unmarshal(endpoints[endpoint], &paramData)

	if paramData[param] == nil {
		return 599, nil, nil, false, "Data for this unique parameter isn't present, returning"
	}

	var data map[string]json.RawMessage
	json.Unmarshal(paramData[param], &data)

	if data["sleep-for-timeout"] != nil {
		var sleep int
		json.Unmarshal(data["sleep-for-timeout"], &sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
	}

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
	paramId := r.URL.Query().Get(os.Getenv("UNIQUE_PARAM_ID"))
	fmt.Printf("Full path is: %s, Unique param is: %s\n", path, paramId)

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
	} else {
		if os.Args[1] == "mode-1" {
			MODE = "mode-1"
		} else if (os.Args[1] == "mode-2") {
			MODE = "mode-2"
		} else {
			fmt.Println("Invalid mode, Exiting!")
			os.Exit(1)
		}
	}

	DATAPATH = "examples"
	MATCHER = "Regex_LCS"
	if len(os.Args) == 3 {
		if os.Args[2] == "samples" {
			DATAPATH = "samples"
		} else if os.Args[2] == "Custom" {
			MATCHER = "Custom"
		}
	} else if len(os.Args) == 4 {
		if os.Args[2] == "samples" {
			DATAPATH = "samples"
		}
		if os.Args[3] == "Custom" {
			MATCHER = "Custom"
		}
	}

	fmt.Printf("Starting HTTP Server on Port: %d, Mode: %s, Datapath: %s, Matcher: %s\n", 3000, MODE, DATAPATH, MATCHER)
	err := http.ListenAndServe(":3000", nil)
        if err != nil {
                panic(err.Error())
        }
}

func main() {
	httpServer()
}
