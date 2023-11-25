package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Card represents the structure of a card with front and back content.
type Card struct {
	ID    string `json:"id"`
	Front string `json:"front"`
	Back  string `json:"back"`
}

var strategies = map[string]func(Card, string) error{
	"md-json": strategyMdJson,
}

var Version = "dev"

func main() {
	dir := flag.String("dir", ".", "directory to scan for markdown files")
	strategy := flag.String("strategy", "md-json", "strategy for processing markdown files")
	outDir := flag.String("out-dir", "out", "output directory for processed files")
	help := flag.Bool("help", false, "print help message")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	// Override flags with .env variables if present.
	loadEnvVariables(dir, strategy, outDir)
	checkArgs(dir, strategy, outDir)
	processMarkdownFiles(strategies[*strategy], *dir, *outDir)
}

// checkArgs checks if the provided arguments are valid.
func checkArgs(dir, strategy, outDir *string) {
	if *dir == "" {
		fmt.Fprintln(os.Stderr, "dir cannot be empty")
		os.Exit(1)
	}
	if *strategy == "" {
		fmt.Fprintln(os.Stderr, "strategy cannot be empty")
		os.Exit(1)
	}
	if _, ok := strategies[*strategy]; !ok {
		fmt.Fprintf(os.Stderr, "invalid strategy: %s\n", *strategy)
		fmt.Fprintln(os.Stderr, "valid strategies:")
		for strategy := range strategies {
			fmt.Fprintf(os.Stderr, "  - %s\n", strategy)
		}
		os.Exit(1)
	}
	if *outDir == "" {
		fmt.Fprintln(os.Stderr, "out-dir cannot be empty")
		os.Exit(1)
	}
}

// loadEnvVariables overrides the provided pointers with values from the .env file if it exists.
func loadEnvVariables(dir, strategy, outDir *string) {
	envFile, err := os.Open(".env")
	if err != nil {
		return
	}
	defer envFile.Close()

	envScanner := bufio.NewScanner(envFile)
	for envScanner.Scan() {
		line := envScanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		switch strings.ToLower(key) {
		case "hash_card_dir":
			*dir = value
		case "hash_card_strategy":
			*strategy = value
		case "hash_card_out_dir":
			*outDir = value
		}
	}
}

// processMarkdownFiles walks through the directory and processes markdown files.
func processMarkdownFiles(
	strategyFn func(Card, string) error,
	dir, outDir string) {

	// Walk through the directory and process markdown files.
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		if d.IsDir() && path == outDir {
			return fs.SkipDir
		}

		processMarkdownFile(strategyFn, path, outDir)

		return nil
	})
}

// Regular expression to match the card line.
var cardRegex = regexp.MustCompile(`#card <!--(\d{4}/\d{2}/\d{2}/[a-zA-Z0-9]+)-->`)

// processMarkdownFile processes a single markdown file.
func processMarkdownFile(
	strategyFn func(Card, string) error,
	path string, outDir string) {
	// Read the content of the markdown file.
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Split the content into lines.
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		// Check if the line matches the card pattern.
		if matches := cardRegex.FindStringSubmatch(line); matches != nil {
			// Extract the card ID.
			cardID := matches[1]

			// Extract the front of the card.
			var front string
			if i > 0 {
				front = lines[i-1]
			} else {
				fmt.Fprintf(os.Stderr, "no front for card %s\n", cardID)
			}

			// Extract the back of the card.
			var back strings.Builder
			for j := i + 1; j < len(lines) && lines[j] != "---"; j++ {
				back.WriteString(lines[j] + "\n")
			}

			// Create the Card object.
			card := Card{
				ID:    cardID,
				Front: front,
				Back:  back.String(),
			}

			err = strategyFn(card, outDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// strategyMdJson processes a single card.
func strategyMdJson(card Card, outDir string) error {
	// Convert the Card to JSON.
	cardJSON, err := json.Marshal(card)
	if err != nil {
		return err
	}

	// Create the output path based on the card ID.
	outputPath := filepath.Join(outDir, card.ID+".md")
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}

	// Write the JSON to the output file.
	if err := os.WriteFile(outputPath, cardJSON, 0644); err != nil {
		return err
	}

	return nil
}
