package main

import (
	"fmt"
	"testing"
)

func TestMain(*testing.T) {
	fmt.Println("testing testing 123...")

	/*

		Things to test

		file name is malitious, e.g. rm*


		input file does not exist
		target root directory does not exist

		editor does not exist

		temp file is created
		temp file is deleted

		indentation is tab/space agnostic
		in a single level, order of file and directories doesn't matter

		file can't be parent of directory or another file

		dry run prints as expected


	*/

}
