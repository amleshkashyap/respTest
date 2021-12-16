package main

import (
	"os"
	"fmt"
	"strconv"
	"net/http"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	_ "github.com/amleshkashyap/respTest/generator"
)

func httpServer() {
	data, err := ioutil.ReadFile("examples/endpoints.json")
        if err != nil {
                fmt.Println("Can't execute without input file")
                os.Exit(1)
        }
        endpoints := make(map[string][]int)
        json.Unmarshal(data, &endpoints)

	for key, val := range endpoints {
		http.HandleFunc(fmt.Sprintf("/%s", key), func(w http.ResponseWriter, r *http.Request) {
			statusCodes := val
			randomIndex := rand.Intn(len(statusCodes))
			code := fmt.Sprintf("%d", statusCodes[randomIndex])

			responseJson, err := ioutil.ReadFile("examples/data.json")
			if err != nil {
				fmt.Printf("Error occurred in reading file: %s\n", err.Error())
				w.WriteHeader(599)
				return
			}

			var someData json.RawMessage
			json.Unmarshal(responseJson, someData)

			var data map[string]json.RawMessage
			json.Unmarshal(responseJson, &data)

			if data[key] == nil {
				fmt.Printf("Data for this endpoint isn't present, returning.\n")
				w.WriteHeader(599)
				return
			}

			var thisData map[string]json.RawMessage
			json.Unmarshal(data[key], &thisData)

			if thisData[code] == nil {
				fmt.Printf("Data for this error code isn't present, finding the first error code with data.\n")
				codeTemp := ""
				for k, _ := range thisData {
					codeTemp = k
					break
				}
				if codeTemp == code {
					fmt.Printf("No error code specific data present for this endpoint, returning.\n")
					w.WriteHeader(599)
					return
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
				fmt.Printf("No actual data present, returning.\n")
				w.WriteHeader(599)
				return
			}

			randomIndex = rand.Intn(len(keyArr))
			fmt.Printf("Returning value for endpoint: %s, code: %s, subcode: %s\n", key, code, keyArr[randomIndex])

			codeVal, err := strconv.Atoi(code)
			if err != nil {
				fmt.Println("Something went wrong")
				w.WriteHeader(599)
				return
			}

			w.WriteHeader(codeVal)
			w.Write(actualData[keyArr[randomIndex]])
			// generator.GenerateStructs("examples/yml_test.yml")
			return
		})
	}

        fmt.Println("Starting HTTP Server on Port: ", 3000)
        err = http.ListenAndServe(":3000", nil)
        if err != nil {
                panic(err.Error())
        }
}

func main() {
	httpServer()
}
