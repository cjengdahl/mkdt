package main

import (
	"bufio"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
)

func createFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func createDirectory(dirname string) error {
	err := os.Mkdir(dirname, 0755) // 0755 is the default permission mode
	return err
}

func main() {

	/*

		todo: input validation
		- basic arg len check
		- input file exists
		- file is not in system root, if it is prompt are you sure?

	*/

	// name := "mkdt"
	// usage := "mkdt [-d] [-r] [-y] [-f]"
	// description := "make directory tree"

	interactiveMode := true
	var inputFilePath string
	flag.StringVar(&inputFilePath, "f", "", "path to input file")

	// rootDirectory := flag.String("r", "", "root directory")
	// dryRun := flag.Bool("d", false, "dry run")
	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	logLevel := new(slog.LevelVar)
	if *verbose {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	if inputFilePath != "" {
		interactiveMode = false
	}

	if interactiveMode {
		// todo: user terminal default editor
		// todo: lead temp file with descriptive comment
		editor := "vim"
		tmpDir := os.TempDir()
		tmpFile, err := os.CreateTemp(tmpDir, "mkdt")
		if err != nil {
			logger.Warn("failed to create temp file", "error", err.Error())
			os.Exit(1)
		}

		path, err := exec.LookPath(editor)
		if err != nil {
			logger.Warn(
				"failed to look up editor path",
				"editor", editor,
				"error", err.Error())
			os.Exit(1)
		}

		cmd := exec.Command(path, tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			logger.Warn(
				"failed to start editor",
				"editor", editor,
				"path", path,
				"error", err.Error())
			os.Exit(1)
		}
		err = cmd.Wait()
		if err != nil {
			logger.Warn(
				"editor encountered error",
				"editor", editor,
				"path", path,
				"error", err.Error())
			os.Exit(1)
		}
		inputFilePath = tmpFile.Name()
	}

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		logger.Warn("error opening input file",
			"path", inputFilePath,
			"error", err.Error())

		os.Exit(1)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	stack := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		// count leading spaces/tabs to determine the indentation level
		indentLevel := 0
		for i := 0; i < len(line) && (line[i] == ' ' || line[i] == '\t'); i++ {
			indentLevel++
		}

		// pop items from the stack until we reach the correct indentation level
		for len(stack) > indentLevel {
			logger.Debug(
				"stack pre-pop",
				"stack", stack,
				"indentLevel", indentLevel,
			)
			stack = stack[:len(stack)-1]
			logger.Debug(
				"stack post-pop",
				"stack", stack,
				"indentLevel", indentLevel,
			)
		}

		// get the directory or file name from the line
		name := strings.TrimSpace(line)

		// join the current stack with the new name to get the full path
		fullPath := path.Join(append(stack, name)...)

		// check if the name contains a dot, skip first character
		if strings.Contains(name[1:], ".") {
			err := createFile(fullPath)
			if err != nil {
				logger.Debug(
					"error creating file",
					"path", fullPath,
					"error", err.Error(),
				)
			} else {
				logger.Debug(
					"created file",
					"path", fullPath,
				)
			}
		} else {
			err := createDirectory(fullPath)
			if err != nil {
				logger.Debug(
					"error creating directory",
					"path", fullPath,
					"error", err.Error(),
				)
			} else {
				logger.Debug(
					"created directory",
					"path", fullPath,
				)
			}
		}

		// push the current name onto the stack
		stack = append(stack, name)
	}

	if err := scanner.Err(); err != nil {
		logger.Debug(
			"error reading input file",
			"error", err.Error(),
		)
	}
}
