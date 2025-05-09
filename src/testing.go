
package main

import "strings"

func testingProgramFromBuf(buf string) Program[MockTerminal] {
	buffers := []Buffer{
		{
			filepath:          "test",
			lines:             strings.Split(buf, "\n"),
			topVisibleLineIdx: 0,
		},
	}

	program := Program[MockTerminal]{
		logger:   getLogger("./logfile.log.txt"),
		state:    ProgramState{},
		term:     MockTerminal{},
		settings: defaultSettings(),
	}

	program.state.buffers = buffers
	initializeState(&program)
	return program
}

const basicBuf =
`package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Terminal interface {
	clearScreen()
	setCursorPosition(x, y int)
	getCursorPosition() (x, y int, err error)
	getSize() (rows, cols int, err error)
	printf(s string, args ...interface{})
}

func (ANSI) getCursorPosition() (x, y int, err error) {
	// Querying the terminal for cursor position
	fmt.Print("\033[6n")

	// Reading the response
	var response []byte
	buf := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read from stdin: %v", err)
		}
		if buf[0] == 'R' {
			break
		}
		response = append(response, buf[0])
	}

	// Parsing the response
	// Response format: "\033[<rows>;<cols>R"
	parts := strings.Split(strings.Trim(string(response), "\033["), ";")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected response format: %s", response)
	}

	rows, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse rows: %v", err)
	}

	cols, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse cols: %v", err)
	}

	return rows, cols, nil
}`
