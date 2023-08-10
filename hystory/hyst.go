//Пн июл 31 16:02:13 MSK 2023
package maklib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var HystDescr = `
About a file "hystory.txt"
The file must be placed in a proces worcing directory and must content hystory records.
A hystory record is a set of lines that
has as first line the head of the record.
A head of a record is a line that matches the pattern of "^<\d\d\d\d\d\d \d\d:\d\d>$",
e. g. "<230728 12:40>" if it is the only content of the line (besides \n).
That is the file's first line must bear the head.
The file may have empty lines and comment lines, that are  lines begin with "#"
A line that is not bare a head may have an arbitrary content
But between a head and the next one mest be at least one info line.
A info line is that that has no blank charackters before, maybe, the charicter "#"
`

type Head struct {
	H    string //A head
	Lbeg int    //A number of a line where a heah of a hystory record is placed
	Lend int    //A number of a last line  of a hystory record
}

//Checks the validity of a Header
//Erros:
//1-Hesder.Valid: bad format;
func (h *Head) Valid() (err error) {
	if h.Lbeg != 0 && h.Lend != 0 && h.H != "" {
		return
	} else {
		err = fmt.Errorf("1$)Hesder.Valid: bad format;%v", h)
		return
	}
	return
}

//Returns the true if all fields are not zero and the h itself is not nil
func (h *Head) IsZero() bool {
	if h != nil && h.Lbeg != 0 && h.Lend != 0 && h.H != "" {
		return false
	} else {
		return true
	}
	return false
}

const hystpattern = "^<\\d\\d\\d\\d\\d\\d\\d\\d:\\d\\d>$"

//Ср авг  9 17:46:00 MSK 2023
//It scans the "hystory.txt" and return heads of hystory records
//Errors
//1 - no file hystory.txt
//2 - regexp.MatchString err
//3 - no heads
//4 - error during scanning
func FindHeads() (heads []*Head, err error) {
	var file *os.File
	var line string
	var lineCount int
	//var matched bool
	//var head string
	var prevHead, currHead *Head //The previous and current Headers

	file, err = os.Open("hystory.txt")
	if err != nil {
		err = fmt.Errorf("1$)hyst.FindHeads: err=%s", err.Error())
		return
	}
	defer file.Close()

	prevHead = new(Head)
	currHead = new(Head)

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
		if _, err = regexp.MatchString(hystpattern, line); err != nil {
			err = fmt.Errorf("2$)hyst.FindHeads:regexp.MatchString err=%s", err.Error())
			return
		} else { //A line matches and there is not error of matching
			if lineCount > 1 && prevHead.IsZero() { //First line must bear the Head
				err = fmt.Errorf("3$)hyst.FindHeads:regexp.MatchString err=%s", err.Error())
				return
			} else {
				currHead.H = line
				currHead.Lbeg = lineCount
				if !prevHead.IsZero() {
					if prevHead.Lbeg < 1 { //It does not be!
						err = fmt.Errorf("4$)hyst.FindHeads:???prevHead.Lbeg, lincont=%d", lineCount)
						return
					}
					currHead.Lend = prevHead.Lbeg - 1
				}
				if lineCount > 1 { //There is not first record
					prevHead = currHead
				}
			}
			return //a case of matchig has been worked out
		} //A line matches and there is not error of matching
		//Further we have deel with a info line
		currHead.Lend = currHead.Lend + 1
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
