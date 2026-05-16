package main

type Config struct {
	Rules map[string][]string `yaml:"rules"`
	Undo  struct {
		LastSortingPath string            `yaml:"last_sorting_path"`
		SortedTo        map[string]string `yaml:"sorted_to"`
	}
}

type Options struct {
	Path   string
	DryRun bool
	Undo   bool
}
