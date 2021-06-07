package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func createConfigFileForTest(t *testing.T, cfg *Config) (string, error) {
	var strCfg string
	if cfg.DBFileName != "" {
		strCfg += fmt.Sprintf("DBFileName: %s\n", cfg.DBFileName)
	}
	if cfg.LogFileName != "" {
		strCfg += fmt.Sprintf("LogFileName: %s\n", cfg.LogFileName)
	}
	if cfg.LogOutputMode != 0 {
		strCfg += fmt.Sprintf("LogOutputMode: %d\n", cfg.LogOutputMode)
	}
	if cfg.LogLevel != "" {
		strCfg += fmt.Sprintf("LogLevel: %s\n", cfg.LogLevel)
	}

	fileName, err := MakeRandomStr(20)
	if err != nil {
		return "", err
	}
	fileName = fileName + ".yaml"
	t.Logf("fileName: %s\n", fileName)
	t.Logf("yaml: \n%s\n", strCfg)
	return fileName, ioutil.WriteFile(fileName, []byte(strCfg), 0777)
}

func checkParam(t *testing.T, createdCfg *Config, expectedCfg *Config) {
	if !reflect.DeepEqual(createdCfg, expectedCfg) {
		t.Fatalf("real: %+v\n expected: %+v\n", *createdCfg, *expectedCfg)
	}
}

func TestNoConfigFile(t *testing.T) {
	// expect that confg having all default params will be created.
	cfg, err := NewConfig("")
	if err != nil {
		t.Fatal(err)
	}
	checkParam(t, cfg, createDefaultConfig())
}

func TestSetPartOfConfig(t *testing.T) {
	// custom config file is created.
	dbFileName := "/opt/tinyurl/customdb.db"
	logOutputMode := 2
	cfg := &Config{
		DBFileName:    dbFileName,
		LogOutputMode: logOutputMode,
	}
	fileName, err := createConfigFileForTest(t, cfg)
	defer os.Remove(fileName)
	if err != nil {
		t.Fatal(err)
	}
	// create pointer of struct unmarshaled above created config file.
	createdCfg, err := NewConfig(fileName)
	if err != nil {
		t.Fatal(err)
	}

	// create pointer of expected struct.
	expectedCfg := createDefaultConfig()
	expectedCfg.DBFileName = dbFileName
	expectedCfg.LogOutputMode = logOutputMode

	t.Logf("config:\n%+v\n", *createdCfg)

	checkParam(t, createdCfg, expectedCfg)
}
