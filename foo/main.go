package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)


func main() {
	// Open the YAML file
    // fmt.Printf("num of args %v\n", os.Args)
    if len(os.Args) == 1 {
        fmt.Println("Supply a file name to parse")
        os.Exit(1)
    }

	file, err := os.Open(os.Args[1]) 
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Decode YAML into Template struct
	var tmpl Template
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&tmpl); err != nil {
		fmt.Println("Error:", err)
		return
	}
    flow := newFrom(tmpl)
    //flow := v1alpha08.Flow{}

    j, err := json.Marshal(&flow)
    if err != nil {
        fmt.Printf("error when marshaling json %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("%v\n", string(j))
    
//    out, err := yaml.Marshal(&flow)
//    if err != nil {
//        fmt.Printf("error when marshaling yaml %v\n", err)
//        os.Exit(1)
//    }
//
//    fmt.Printf("%v\n", string(out))

}
