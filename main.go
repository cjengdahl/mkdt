package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type treeNode struct {
	isDir    bool
	depth    int
	name     string
	children []*treeNode
}

// create an empty file
func createFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// create an empty directory
func createDirectory(dirname string) error {
	return os.Mkdir(dirname, 0755)
}

// create/write input file with editor
func makeInputFileWithEditor() string {

	editor := os.Getenv("EDITOR")
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "mkdt")
	if err != nil {
		fmt.Printf("failed to create temp file: %s\n", err.Error())
		os.Exit(1)
	}

	path, err := exec.LookPath(editor)
	if err != nil {
		fmt.Printf("failed to look up editor path: %s\n", err.Error())
		os.Exit(1)
	}

	cmd := exec.Command(path, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Printf("failed start editor: %s\n", err.Error())
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("editor encountered: %s\n", err.Error())
		os.Exit(1)
	}
	return tmpFile.Name()
}

// build in-memory file/directory tree
func buildTree(rootDirPath, inputFilePath string) treeNode {

	root := treeNode{
		isDir:    true,
		depth:    0,
		name:     rootDirPath,
		children: []*treeNode{},
	}
	stack := []*treeNode{&root}
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Printf("error opening input file: %s\n", err.Error())
		os.Exit(1)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {

		line := scanner.Text()
		name := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		isDir := false
		if !strings.Contains(name[1:], ".") {
			isDir = true
		}
		indentLevel := 0
		for i := 0; i < len(line) && (line[i] == ' ' || line[i] == '\t'); i++ {
			indentLevel++
		}

		newNode := treeNode{
			isDir:    isDir,
			depth:    indentLevel + 1,
			name:     name,
			children: []*treeNode{},
		}

		for {

			// sibling of head of stack
			if newNode.depth == stack[len(stack)-1].depth {
				stack[len(stack)-2].children = append(stack[len(stack)-2].children, &newNode)
				stack[len(stack)-1] = &newNode // overwrite sibling in stack, sibling has no more children
				break
			}

			// child of head of stack
			if newNode.depth-stack[len(stack)-1].depth >= 1 {
				if !stack[len(stack)-1].isDir {
					fmt.Printf("invalid tree: %s is a file, but has %s as a child file/directory\n",
						stack[len(stack)-1].name,
						newNode.name,
					)
					os.Exit(1)
				}
				stack[len(stack)-1].children = append(stack[len(stack)-1].children, &newNode)
				stack = append(stack, &newNode)
				break
			}

			stack = stack[:len(stack)-1]

		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading input file: %s\n", err.Error())
		os.Exit(1)
	}

	return root

}

// create file/directory tree by traversing the in-memory tree
func createTree(rootDirPath, parentPath string, node *treeNode) {
	if node.isDir {
		dirPath := path.Join(parentPath, node.name)
		if dirPath != rootDirPath {
			err := createDirectory(path.Join(parentPath, node.name))
			if err != nil {
				fmt.Printf("aborting, failed create to directory %s:\n", err.Error())
				os.Exit(1)
			}
		}
	} else {
		err := createFile(path.Join(parentPath, node.name))
		if err != nil {
			fmt.Printf("aborting, failed to create file %s:\n", err.Error())
			os.Exit(1)
		}
	}
	for _, child := range node.children {
		createTree(rootDirPath, path.Join(parentPath, node.name), child)
	}
}

func main() {

	usage := func() {
		fmt.Println("description:\n\tmake a directory tree from a text file (interactivley by default)")
		fmt.Println("\tdelimited by newline (new directory/file) and space (level in the tree)")
		fmt.Println("usage:\n\tmkdt [-fr]")
		fmt.Println("example:\n\tmkdir")
		fmt.Println("\ta\n\t 1.txt\n\tb\n\t c\n\t  2.txt\n\t3.txt")
	}

	flag.Usage = usage

	var inputFilePath string
	flag.StringVar(&inputFilePath, "f", "", "path to input file")

	var rootDirPath string
	flag.StringVar(&rootDirPath, "r", ".", "root directory to create tree")

	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Println("mkdt takes no arguments, flags are optional")
		flag.Usage()
		os.Exit(1)
	}

	if inputFilePath == "" {
		inputFilePath = makeInputFileWithEditor()
	}

	root := buildTree(rootDirPath, inputFilePath)
	createTree(rootDirPath, "", &root)

}
