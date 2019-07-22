package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gesellix/couchdb-prometheus-exporter/lib"
)

func main() {
	updateSeq := os.Args[1]
	//fmt.Printf("%s\n", updateSeq)

	decoded, err := lib.DecodeUpdateSeq(updateSeq)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%v\n", decoded)

	decodedJson, _ := json.Marshal(decoded)
	fmt.Printf("%s\n", decodedJson)
}
