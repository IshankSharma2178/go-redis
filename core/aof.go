package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/IshankSharma2178/go-redis/internals/config"
)

// TODO: Support Expiration
// TODO: Support non-kv data structures
// TODO: Support sync write
func dumpKey(fp *os.File, key string, obj *Obj) {
	cmd := fmt.Sprintf("SET %s %s", key, obj.Value)
	tokens := strings.Split(cmd, " ")
	fp.Write(Encode(tokens, false))
}

// TODO: To to new and switch
func DumpAllAOF() {
	if err := os.MkdirAll(filepath.Dir(config.Cfg.AOFFile), 0o755); err != nil {
		fmt.Print("error", err)
		return
	}

	fp, err := os.OpenFile(config.Cfg.AOFFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Print("error", err)
		return
	}
	defer fp.Close()
	log.Println("rewriting AOF file at", config.Cfg.AOFFile)
	for k, obj := range store {
		dumpKey(fp, k, obj)
	}
	log.Println("AOF file rewrite complete")
}
