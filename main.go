package main

import (
	"fmt"
	"os"
)

func main() {
	cfg, err := NewConfig("")
	defer func() {
		if logger != nil {
			Errorf("MainError: %v\n", err)
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "MainError: %v\n", err)
		}
	}()
	logger, err = SetupLogger(cfg.LogFileName, cfg.LogOutputMode, cfg.LogLevel)
	db, err := ConnectDB(cfg.DBFileName)
	err = StartTinyURLServer(cfg, db)
}
