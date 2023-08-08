package core

import (
	"fmt"
	"log"
	"memkv/config"
	"os"
	"strings"
)

func DumpAllAOF() {
	f, err := os.OpenFile(config.AOFFileName, os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Print("error when creating AOF file: ", err)
		return
	}
	for k, o := range store {
		cmd := fmt.Sprintf("SET %s %s", k, o.Value)
		tokens := strings.Split(cmd, " ")
		f.Write(Encode(tokens, false))
	}
	log.Println("AOF file rewrite done")
}
