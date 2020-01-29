package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	slash   byte = 47
	space   byte = 32
	star    byte = 42
	newline byte = 10
)

func readFile(fPath string) (string, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	b, err = strip(b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func strip(b []byte) ([]byte, error) {
	//TODO -> remove multiline comments
	singleLineComment := false
	var b1 []byte

	for i := 0; i < len(b); i++ {
		if (b[i] == slash) && (i+1 < len(b)) {
			if b[i+1] == slash {
				singleLineComment = true
				i++
			}
		} else if b[i] == newline {
			singleLineComment = false
		} else if b[i] == space {
			continue
		} else if !singleLineComment {
			b1 = append(b1, b[i])
		}
	}

	return b1, nil
}

func capFirst(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func rec(name string, m map[string]interface{}, d int) (string, error) {
	var prefix string
	if name == "" {
		prefix = `type NameHere struct `
	} else {
		prefix = capFirst(name) + `struct `
	}

	for key := range m {
		switch m[key].(type) {
		case map[string]interface{}:
			fmt.Println("map ->", key)
		default:
			fmt.Println("normal field:", key)
		}
	}

	return prefix, nil
}

func extractStruct(s string) (string, error) {
	//TODO
	var res map[string]interface{}

	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		fmt.Println("Something wrong with your JSON")
		return "", err
	}

	recRet, err := rec("", res, 0)
	if err != nil {
		return "", err
	}

	fmt.Println(recRet)

	return s, nil
}

func main() {
	fPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fPath = fPath + "/temp.txt"

	f, err := os.Create(fPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	_, err = f.WriteString("// ------------- Welcome to Elle-Station ---------------" +
		"\n// Paste your JSON below and save & exit the text editor" +
		"\n// -----------------------------------------------------")
	f.Close()

	defer os.Remove(fPath)

	editor := flag.String("editor", "nano", "Set the editor that use is presented with to paste their JSON")
	flag.Parse()

	cmd := exec.Command(*editor, fPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Println("Could not open editor:", *editor)
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error while editing. Error:", err)
		os.Exit(1)
	}

	body, err := readFile(fPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	finalStruct, err := extractStruct(body)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(finalStruct)

	// Ask user if they want it saved to a go file
}
