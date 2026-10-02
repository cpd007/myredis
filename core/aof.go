package core

import (
	"fmt"
	"log"
	"os"

	"github.com/cpd007/myredis/config"
)

func DumpAllAOF() {

	fp, err := os.OpenFile(config.AOFFile, os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Print("error in opening file", err)
		return
	}

	log.Println("rewriting AOF file at", config.AOFFile)

	for k, obj := range store {
		dumpKey(fp, k, obj)
	}

	log.Println("AOF rewrite complete")
}

func dumpKey(fp *os.File, k string, obj *Obj) {

	tokens := []string{"SET", k, fmt.Sprintf("%s", obj.Value)}
	encodedData, err := Encode(tokens, false)
	if err != nil {
		log.Println("error in encoding keys for dump file", err, k, obj)
		return
	}
	fp.Write(encodedData)
}
