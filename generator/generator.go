package generator

import (
	"encoding/json"
        "fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
)

var AllTypes = [...]string {"int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"float32", "float64",
	"complex64", "complex128",
	"byte", "rune", "string", "bool"}

func GenerateStructs(filename string) {
	fmt.Println(filename)
	fmt.Println()

	file, _ := ioutil.ReadFile(filename)
        obj := make(map[string]map[string]interface{})

	// []string doesn't work in yml files, use "[]string"
	if strings.Contains(filename, ".yml") {
		yaml.Unmarshal(file, &obj)
	} else if strings.Contains(filename, ".json") {
		json.Unmarshal(file, &obj)
	} else {
		fmt.Printf("Unexpected Source File, Exiting Unexpectedly.")
		os.Exit(1)
	}

	var AllStructs []string
	for k, _ := range obj {
		AllStructs = append(AllStructs, k)
	}

	for k, v := range obj {
		fmt.Printf("type %s struct {\n", k)
		for key, value := range v {
			found := false
			for _, t := range AllTypes {
				if ( (t == value) || ((fmt.Sprintf("[]%s", t)) == value) ) {
					fmt.Printf("\t%s %s\n", key, value)
					found = true
					break
				}
			}

			if found == false {
				for _, t := range AllStructs {
					if ( (t == value) || ((fmt.Sprintf("[]%s", t)) == value) ) {
						fmt.Printf("\t%s %s\n", key, value)
						found = true
						break
					}
				}
			}
			if found == false {
				fmt.Println("Invalid or Unsupported Datatype, Exiting Unexpectedly.")
				os.Exit(1)
			}
		}
		fmt.Printf("}\n\n")
	}
}
