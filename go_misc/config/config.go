package config

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/TrungACZNE/go_misc/cliparser"

	"log"
)

var (
	intVars   = make(map[string]int64)
	strVars   = make(map[string]string)
	floatVars = make(map[string]float64)
	boolVars  = make(map[string]bool)

	intLock   = &sync.Mutex{}
	strLock   = &sync.Mutex{}
	floatLock = &sync.Mutex{}
	boolLock  = &sync.Mutex{}
)

func IsUnsetOrFalse(varname string) bool {
	boolLock.Lock()
	defer boolLock.Unlock()
	if val, ok := boolVars[varname]; !ok || val == false {
		return true
	}
	return false
}

// missingVar convenient method to exit program & report missing var
func missingVar(varname string) {
	log.Fatalf(`Attempted to get unset variable "%s"`, varname)
}

// MustGetInt returns an int64 config variable atomically. Fatal if missing.
func MustGetInt(varname string) int64 {
	intLock.Lock()
	defer intLock.Unlock()
	if val, ok := intVars[varname]; ok {
		return val
	}
	missingVar(varname)
	return -1
}

// GetFloat returns a float64 config variable atomically. Fatal if missing.
func MustGetFloat(varname string) float64 {
	floatLock.Lock()
	defer floatLock.Unlock()
	if val, ok := floatVars[varname]; ok {
		return val
	}
	missingVar(varname)
	return -1.0
}

// MustGetString returns a string config variable atomically. Fatal if missing.
func MustGetString(varname string) string {
	strLock.Lock()
	defer strLock.Unlock()
	if val, ok := strVars[varname]; ok {
		return val
	}
	missingVar(varname)
	return ""
}

// MustGetBool returns a bool config variable atomically. Fatal if missing.
func MustGetBool(varname string) bool {
	boolLock.Lock()
	defer boolLock.Unlock()
	if val, ok := boolVars[varname]; ok {
		return val
	}
	missingVar(varname)
	return false
}

// GetInt returns an int64 config variable atomically.
func GetInt(varname string) (int64, error) {
	intLock.Lock()
	defer intLock.Unlock()
	if val, ok := intVars[varname]; ok {
		return val, nil
	}
	return -1, fmt.Errorf("Unset var")
}

// GetFloat returns a float64 config variable atomically.
func GetFloat(varname string) (float64, error) {
	floatLock.Lock()
	defer floatLock.Unlock()
	if val, ok := floatVars[varname]; ok {
		return val, nil
	}
	return -1.0, fmt.Errorf("Unset var")
}

// GetStr returns a string config variable atomically.
func GetStr(varname string) (string, error) {
	strLock.Lock()
	defer strLock.Unlock()
	if val, ok := strVars[varname]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unset var")
}

// GetBool returns a bool config variable atomically.
func GetBool(varname string) (bool, error) {
	boolLock.Lock()
	defer boolLock.Unlock()
	if val, ok := boolVars[varname]; ok {
		return val, nil
	}
	return false, fmt.Errorf("Unset var")
}

// SetInt assigns or overwrites an int64 value to a config variable atomically
func SetInt(varname string, value int64) {
	intLock.Lock()
	defer intLock.Unlock()

	intVars[varname] = value
}

// SetFloat assigns or overwrites a float64 value to a config variable atomically
func SetFloat(varname string, value float64) {
	floatLock.Lock()
	defer floatLock.Unlock()

	floatVars[varname] = value
}

// SetStr assigns or overwrites a string value to a config variable atomically
func SetStr(varname string, value string) {
	strLock.Lock()
	defer strLock.Unlock()

	strVars[varname] = value
}

// SetBool assigns or overwrites a string value to a config variable atomically
func SetBool(varname string, value bool) {
	boolLock.Lock()
	defer boolLock.Unlock()

	boolVars[varname] = value
}

func Set(varname string, vartype string, varvalue string) error {
	switch vartype {
	case "I":
		val, err := strconv.ParseInt(varvalue, 10, 64)
		if err != nil {
			return err
		} else {
			SetInt(varname, val)
		}
	case "F":
		val, err := strconv.ParseFloat(varvalue, 64)
		if err != nil {
			return err
		} else {
			SetFloat(varname, val)
		}
	case "S":
		SetStr(varname, varvalue)
	case "B":
		val, err := strconv.ParseBool(varvalue)
		if err != nil {
			return err
		} else {
			SetBool(varname, val)
		}
	default:
		return fmt.Errorf("Type must be B, I, F or S")
	}
	return nil
}

func LoadFromFile(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	text := string(bytes)
	for lnum, line := range strings.Split(text, "\n") {
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		params, err := cliparser.ParseCommandless(line)
		if len(params) != 3 {
			log.Printf("Line number %v needs to have 3 parameters: name - type - value\n", lnum)
		} else if err != nil {
			log.Printf("Error '%v':\n[%s:%v] \n", err, filename, lnum)
		} else {
			err := Set(params[0], params[1], params[2])
			if err != nil {
				log.Fatalf("Error '%v':\n[%s:%v] \n", err, filename, lnum)
			}
		}

	}
	return nil
}
