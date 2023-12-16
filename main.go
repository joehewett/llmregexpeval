package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose output")
	help       = flag.Bool("help", false, "help")
	promptFile = flag.String("prompt", "", "The file to read the prompt from. Prompt must be in OpenAI Chat format and contain {{ .task }} in the text.")
	rulesFile  = flag.String("rules", "", "A JSON file containing the regex rules to run over LLM outputs")
	taskFile   = flag.String("tasks", "", "A JSON file containing tasks to substitute into the master prompt")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage: go run main.go (-d <directory>|-f <file>) -j <json file>")
		return
	}

	prompt, err := loadPrompt()
	if err != nil {
		fmt.Printf("Failed to load JSON: %s\n", err)
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Failed to load data: %s\n", err)
		return
	}

	startTime := time.Now()
	defer func() {
		if *verbose {
			fmt.Printf("Total time taken: %s\n", time.Since(startTime))
		}
	}()

	ch := make(chan string)
	for _, task := range tasks {
		go fill(prompt, task, ch)
	}

	for range tasks {
		fmt.Printf(<-ch)
	}

}

func fill(file string, jsonStr string, ch chan string) error {
	startTime := time.Now()

	if _, err := os.Stat(file); err != nil {
		return fmt.Errorf("file %s does not exist", file)
	}

	if *verbose {
		fmt.Println("Reading data from file")
	}

	dataFile, err := os.Open(file)
	if err != nil {
		return err
	}

	defer dataFile.Close()

	bytes, err := io.ReadAll(dataFile)

	if err != nil {
		fmt.Println(err)
		return err
	}

	data := string(bytes)
	if *verbose {
		fmt.Println("Reading JSON from file")
	}

	result, err := requestFill(jsonStr, data)
	if err != nil {
		fmt.Printf("Failed to request fill for file %s: %s\n", file, err)
		return nil
	}

	ch <- result
	if *verbose {
		fmt.Printf("Time taken for file %s: %s\n", file, time.Since(startTime))
	}

	return nil
}

func loadPrompt() (string, error) {
	if _, err := os.Stat(*promptFile); err != nil {
		return "", fmt.Errorf("failed to read prompt file")
	}

	if *verbose {
		fmt.Printf("Reading prompt JSON from file %s\n", *promptFile)
	}

	handle, err := os.Open(*promptFile)
	if err != nil {
		return "", err
	}

	defer handle.Close()

	bytes, err := io.ReadAll(handle)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bytes, new((map[string]interface{})))
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal JSON skeleton, please check your JSON is valid and try again: %s\n", err)
	}

	promptJSON := string(bytes)
	if *verbose {
		fmt.Printf("JSON: %s\n", promptJSON)
	}

	if promptJSON == "" {
		return "", fmt.Errorf("please provide the JSON containing your prompt in OpenAI format")
	}

	return promptJSON, nil
}

func loadTasks() ([]string, error) {
	var taskStrings []string

	if *taskFile == "" {
		return nil, fmt.Errorf("you need to specify a task file containing an array of strings, each representing an LLM task")
	}

	if *verbose {
		fmt.Printf("Reading data from directory %s\n", *taskFile)
	}

	handle, err := os.Open(*taskFile)
	if err != nil {
		return nil, fmt.Errorf("could not open task file")
	}

	bytes, err := io.ReadAll(handle)
	if err != nil {
		return nil, fmt.Errorf("could not read task file")
	}

	err = json.Unmarshal(bytes, &taskStrings)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON skeleton, please check your JSON is valid and try again: %s\n", err)
	}

	return taskStrings, nil
}
