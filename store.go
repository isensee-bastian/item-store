package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	DATA_FILE     = "items"
	ACTION_ADD    = "add"
	ACTION_REMOVE = "remove"
	ACTION_LIST   = "list"
)

func main() {
	if len(os.Args) < 2 {
		exit("Expected action argument.")
	}

	action := os.Args[1]

	switch action {
	case ACTION_ADD:
		items := os.Args[2:]
		if len(items) == 0 {
			exit("Expected at least one item name to add.")
		}
		add(items)
	case ACTION_REMOVE:
		items := os.Args[2:]
		if len(items) == 0 {
			exit("Expected at least one item name to remove.")
		}
		remove(items)
	case ACTION_LIST:
		list()
	default:
		// Show help info when no other case applied.
		fmt.Println("Valid actions are:", ACTION_ADD, ACTION_REMOVE, ACTION_LIST)
	}
}

func remove(items []string) {
	content, err := os.ReadFile(DATA_FILE)
	check(err)

	// Trim for removing empty line at the end.
	// Split for converting lines into slice elements.
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	remove := make(map[string]bool)
	for _, item := range items {
		remove[item] = true
	}

	file, err := os.Create(DATA_FILE)
	check(err)
	defer file.Close()

	// An alternative to recreating the file with os.Create is to open the file
	// in write mode and truncate its content.
	// file, err := os.OpenFile(DATA_FILE, os.O_CREATE|os.O_WRONLY, 0644)
	// check(err)
	// defer file.Close()
	// Empty the file.
	// err = file.Truncate(0)
	// check(err)

	for index := 0; index < len(lines); index++ {
		// If line item is not contained in map created from user specified items
		// it shall be kept (otherwise removed).
		item := lines[index]
		if !remove[item] {
			file.WriteString(item + "\n")
		}
	}
}

func add(items []string) {
	file, err := os.OpenFile(DATA_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, item := range items {
		_, err = writer.WriteString(item + "\n")
		check(err)
	}

	// With buffered writing we need to remember calling flush
	// for takeover of in memory data to hard drive file.
	err = writer.Flush()
	check(err)
}

func list() {
	// First approach reads all file content into memory
	// which is fine for smaller files. Note that the file has to exist, though.
	content, err := os.ReadFile(DATA_FILE)
	check(err)
	fmt.Println(strings.TrimSpace(string(content)))

	// Second approach reads file line by line which is better
	// for larger files.
	/*
		file, err := os.OpenFile(DATA_FILE, os.O_CREATE|os.O_RDONLY, 0644)
		check(err)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		check(scanner.Err())
	*/
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func check(err error) {
	if err != nil {
		exit(err.Error())
	}
}
