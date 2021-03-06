package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	vaultDisk "github.com/hashicorp/vault/builtin/token/disk"
)

func mainToken(session string, args []string) int {
	if len(args) == 0 {
		return 1
	}

	path := fmt.Sprintf("tmp/sessions/%s", session)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return 2
	}

	log.Printf("[DEBUG] token: %v %s", args, path)
	helper := &vaultDisk.Command{Path: path}
	return helper.Run(args)
}
