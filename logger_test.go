/* todo: create test
- output log file
- per log level
*/
package main

import (
	"os"
	"strings"
	"testing"
)

// intersept and get input of stderr by logger.
func getStderr(t *testing.T, fn func()) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}
	stderr := os.Stderr
	defer func() {
		os.Stderr = stderr
	}()
	os.Stderr = writer
	fn()
	writer.Close()
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func TestLoggerWrite(t *testing.T) {
	fileName, err := MakeRandomStr(15)
	if err != nil {
		t.Fatal(err)
	}
	fileName = "/tmp/" + fileName + ".log"
	if logger, err = SetupLogger(fileName, 0, "debug"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName)

	if logger == nil {
		t.Fatal("logger is nil")
	}

	exp := "Hello world."
	str, err := getStderr(t, func() {
		logger.Write([]byte(exp))
	})
	if err != nil {
		t.Fatal(err)
	}
	if str != exp {
		t.Fatalf("\nreal: %s\nexpected: %s\n", str, exp)
	}
}

func TestLogFormatCheck(t *testing.T) {
	fileName, err := MakeRandomStr(15)
	if err != nil {
		t.Fatal(err)
	}
	fileName = "/tmp/" + fileName + ".log"
	if logger, err = SetupLogger(fileName, 0, "debug"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName)

	info := "[info] info test"
	str, err := getStderr(t, func() {
		Infof("info %s", "test")
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(str, info) == -1 {
		t.Fatalf("\nreal: %s\nexpected: %s\n", str, info)
	}

	warn := "[warn] warn yeees"
	str, err = getStderr(t, func() {
		Warnf("warn %s", "yeees")
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(str, warn) == -1 {
		t.Fatalf("\nreal: %s\nexpected: %s\n", str, warn)
	}

	erro := "[error] error hoooo"
	str, err = getStderr(t, func() {
		Errorf("error %s", "hoooo")
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(str, erro) == -1 {
		t.Fatalf("\nreal: %s\nexpected: %s\n", str, erro)
	}

	debug := "[debug] debug boooo"
	str, err = getStderr(t, func() {
		Debugf("debug %s", "boooo")
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(str, debug) == -1 {
		t.Fatalf("\nreal: %s\nexpected: %s\n", str, debug)
	}
}
