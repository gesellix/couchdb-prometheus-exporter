package main

import (
	"encoding/json"
	"fmt"
	"github.com/gesellix/couchdb-prometheus-exporter/lib"
	"os"
)

func main() {
	updateSeq := os.Args[1]
	fmt.Printf("** updateSeq:\n%s\n", updateSeq)

	decoded, err := lib.DecodeUpdateSeq(updateSeq)
	if err != nil {
		panic(err)
	}
	decodedJson, _ := json.Marshal(decoded)
	fmt.Printf("** decoded:\n%v\n** json:\n%s\n", decoded, decodedJson)
}
