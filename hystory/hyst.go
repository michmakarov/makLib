//Пн июл 31 16:02:13 MSK 2023
package maklib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const hystpattern = "^<\\d\\d\\d\\d\\d\\d\\d\\d:\\d\\d>$"

var HystDescr = `
In the process it look for file hystory.txt.
The file may content hystory records.
A hystory record is a setof lines that
has as first line the head of the record.
A head of a record os a line that matches the pattern of "^<\d\d\d\d\d\d \d\d:\d\d>$",
e. g. "<230728 12:40>" if it is the only content of the line (besides \n).
`

//Ср авг  9 17:46:00 MSK 2023
//It scans the "hystory.txt" and return heads of hystory records
//Errors
//1 - no file hystory.txt
//2 - regexp.MatchString err
//3 - no heads
//4 - error during scanning
func FindHeads() (heads []string, err error) {
	var file *os.File
	var line string
	var lineCount int
	var matched bool
	file, err = os.Open("hystory.txt")
	if err != nil {
		err = fmt.Errorf("1$)hyst.FindHeads: err=%s", err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		lineCount++
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == byte('#') {
			continue
		}
		if matched, err = regexp.MatchString(hystpattern, line); err != nil {
			err = fmt.Errorf("2$)hyst.FindHeads:regexp.MatchString err=%s", err.Error())
			return
		} else {
			if matched {
				heads = append(heads, line)
			} else {
				continue
			}
		}

	} //for scanner.Scan()
	if len(heads) == 0 {
		err = fmt.Errorf("3$)hyst.FindHeads:no heads, lines=%d", lineCount)
		return
	}

	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("4$)hyst.FindHeads:There is errr while scanning;err=%s;lineCount=%d", err.Error(), lineCount)
		return
	}
	return
}
