package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const CONFIG_NAME string = "config.yaml"

var sourceDirMap = make(map[string][]string)
var extDirMap = make(map[string]string)
var lastSortingPath string
var pathsChan = make(chan map[string]string)
var lastlySortedTo map[string]string

func loadConfig(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Config not found!!! Creating default config...")
			if err = createDefaultConfig(CONFIG_NAME); err != nil {
				return err
			}
			fmt.Println("Created!")
			data, err = os.ReadFile(name)
			if err != nil {
				return err
			}
		}
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	sourceDirMap = cfg.Rules
	for dir, exts := range cfg.Rules {
		for _, ext := range exts {
			extDirMap[ext] = dir
		}
	}

	lastSortingPath = cfg.Undo.LastSortingPath
	lastlySortedTo = cfg.Undo.SortedTo

	return nil
}

func createDefaultConfig(name string) error {
	var cfg Config
	cfg.Rules = map[string][]string{
		"Documents": {".txt", ".md", ".log"},
		"Pictures":  {".jpg", ".png", ".gif", ".svg"},
		"Office":    {".pdf", ".docx", ".xlsx", ".pptx"},
		"Media":     {".mp3", ".wav", ".mp4", ".mkv"},
		"Archives":  {".zip", ".tar.gz", ".rar", ".7z"},
		"Web":       {".html", ".css", ".js", ".json"},
		"Code":      {".sh", ".py", ".cpp", ".go"},
		"Config":    {".conf", ".ini", ".yaml", ".xml"},
	}
	return saveConfigByStruct(cfg, name)
}

func saveConfigByStruct(cfg Config, name string) error {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err = os.WriteFile(name, out, 0644); err != nil {
		return err
	}
	return nil
}

func saveConfig(path, name string, undoPaths map[string]string) error {
	var cfg Config
	cfg.Rules = sourceDirMap
	cfg.Undo.LastSortingPath = path
	cfg.Undo.SortedTo = undoPaths

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err = os.WriteFile(name, out, 0644); err != nil {
		return err
	}
	return nil
}

func parseOptions() *Options {
	opts := &Options{}

	flag.StringVar(&opts.Path, "path", ".", "Path to directory for sorting")
	flag.BoolVar(&opts.DryRun, "dry-run", false, "Show changes without moving files")
	flag.BoolVar(&opts.Undo, "undo", false, "Revert last sorting action")

	flag.Parse()

	return opts
}

func cancelSorting(paths map[string]string) error {
	fmt.Println("Reverting last changes...")
	var counter int
	for new, old := range paths {
		if err := os.Rename(new, old); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		fmt.Printf("%s ==> %s\n", new, old)
		counter++
	}
	if counter == 0 {
		fmt.Println("Nothing to revert!")
	}
	return nil
}

func confirmAction(path string) {
	var userApproved string
	fmt.Printf("Are you sure you want to sort files in %s? [y/N]: ", path)
	fmt.Scanln(&userApproved)
	userApproved = strings.ToLower(userApproved)

	switch userApproved {
	case "y":
	case "n":
		os.Exit(0)
	default:
		os.Exit(0)
	}
}

func sortFilesToDirsByExt(path string, files []os.DirEntry, ch chan<- map[string]string, isDryRun bool) {
	var error, success int
	for _, file := range files {
		filename := file.Name()
		if file.IsDir() {
			fmt.Printf("%s is a directory, skipping...\n", filename)
			continue
		}
		ext := filepath.Ext(filename)
		lowerExt := strings.ToLower(ext)
		dir, ok := extDirMap[lowerExt]
		if !ok {
			extDirMap[lowerExt] = "Other"
			dir = "Other"
		}

		pathToNewDir := filepath.Join(path, dir)
		oldPathToFile := filepath.Join(path, filename)
		newPathToFile := filepath.Join(pathToNewDir, filename)

		i := 1
		originalName := filename
		for {
			_, err := os.Stat(newPathToFile)
			if os.IsNotExist(err) {
				break
			}

			// if in new directory already exists file with the same name:
			// create filename with unused index
			ext := filepath.Ext(originalName)
			base := strings.TrimSuffix(originalName, ext)
			filename = fmt.Sprintf("%s_%d%s", base, i, ext)
			newPathToFile = filepath.Join(pathToNewDir, filename)

			i++
		}

		if !isDryRun {
			if err := os.MkdirAll(pathToNewDir, 0o755); err != nil {
				fmt.Printf("Failed to create directory %s: %v", pathToNewDir, err)
				error++
			}
			if err := os.Rename(oldPathToFile, newPathToFile); err != nil {
				fmt.Printf("Failed to move file to directory %s: %v", pathToNewDir, err)
				error++
			}
		}

		ch <- map[string]string{newPathToFile: oldPathToFile}

		fmt.Printf("%s ==> %s\n", filename, pathToNewDir)
		success++
	}
	fmt.Printf("--- FINISH ---\nFinished sorting successfully! Totally sorted %d files, errors: %d\n", success, error)
}

func main() {
	var history = make(map[string]string)
	opts := parseOptions()

	if err := loadConfig(CONFIG_NAME); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	files, err := os.ReadDir(opts.Path)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Undo {
		if err := cancelSorting(lastlySortedTo); err != nil {
			log.Fatal(err)
		}
		return
	}

	if !opts.DryRun {
		confirmAction(opts.Path)
	} else {
		fmt.Println("--- DRY RUN MODE: No files will be moved ---")
	}

	go func() {
		for move := range pathsChan {
			for k, v := range move {
				history[k] = v
			}
		}
	}()

	sortFilesToDirsByExt(opts.Path, files, pathsChan, opts.DryRun)

	close(pathsChan)

	saveConfig(opts.Path, CONFIG_NAME, history)
}
