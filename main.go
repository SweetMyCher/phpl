package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func lint(file string) {
	stdout, err := exec.Command("php", "-l", file).Output()
	if err != nil {
		message := string(stdout)
		message = strings.TrimSpace(message)
		fmt.Println(message)
	}
}

func main() {
	err := exec.Command("which", "php").Run()
	if err != nil {
		log.Fatal("command not found: php")
	}

	files := make([]string, 0)

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		path = strings.TrimSpace(path)
		if path == "" || strings.HasPrefix(path, "vendor") || strings.HasPrefix(path, "node_modules") {
			return err
		}
		if strings.HasSuffix(path, ".php") {
			files = append(files, path)
		}
		return err
	})

	jobs := make(chan string, len(files))
	results := make(chan bool, len(files))

	for w := 1; w <= 20; w++ {
		go worker(w, jobs, results)
	}

	for _, file := range files {
		jobs <- file
	}

	close(jobs)

	for a := 1; a <= len(files); a++ {
		<-results
	}
}

func worker(id int, jobs <-chan string, results chan<- bool) {
	for j := range jobs {
		lint(j)
		results <- true
	}
}
