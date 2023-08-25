//Пн июл 31 16:02:13 MSK 2023
package maklib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

var HystDescr = `
About a file "hystory.txt"
It is a text file that tells about significant (from the developer's point of view) stages of develoring some app.
The file must be placed in a proces worcing directory and must content hystory records.
A hystory record is a set of lines that
has as first line the head of the record.
A head of a record is a line that matches the pattern of "^<\d\d\d\d\d\d \d\d:\d\d>$",
e. g. "<230728 12:40>" if it is the only content of the line (besides \n).
That is the file's first line must bear the head.
The file may have empty lines and comment lines, that are  lines begin with "#"
A line that is not bare a head may have an arbitrary content and is named info line.
A info line is one that is not empty or comment line
Between a head and the next one must be at least one info line.
`

type Head struct {
	H    string //A head
	Lbeg int    //A number of a line where a heah of a hystory record is placed
	Lend int    //A number of a last line  of a hystory record
}

//returns err==nil if all fielgs are zero; errors:
//1- 1$)head.IzZero: The Head is nil
//2- 1$)head.IzZero: there is a zero field
func (h *Head) IsZero() (err error) {
	if h == nil {
		err = fmt.Errorf("1$)head.IzZero: The Head is nil")
		return
	}
	if h.Lbeg != 0 || h.Lend != 0 || h.H != "" {
		err = fmt.Errorf("1$)head.IzZero: there is a zero field")
		return
	}
	return
}

// error:
//1 - 1$)head.SetZero: The Head is nil"
func (h *Head) SetZero() (err error) {
	if h == nil {
		err = fmt.Errorf("1$)head.SetZero: The Head is nil")
		return
	}
	h.H = ""
	h.Lbeg = 0
	h.Lend = 0
	return
}

//returns err==nil if only Lbeg and  H is not zero:
//1- 1$)head.IzZero: The Head is nil
//2- 2$)head.Initialized: Lend is not zero
//3- 3$)head.Initialized: 3$)head.Initialized: Lbeg or  h.H  is zero
func (h *Head) Initialized() (err error) {
	if h == nil {
		err = fmt.Errorf("1$)head.Initialized: The Head is nil")
		return
	}
	if h.Lend != 0 {
		err = fmt.Errorf("2$)head.Initialized: Lend is not zero")
		return
	}
	if h.Lbeg == 0 || h.H == "" {
		err = fmt.Errorf("3$)head.Initialized: Lbeg or  h.H  is zero")
		return
	}
	return
}

const hystpattern = "^<\\d\\d\\d\\d\\d\\d\\d\\d:\\d\\d>$"

//Ср авг  9 17:46:00 MSK 2023
//Пн авг 14 08:19:16 MSK 2023 That is third assult
//An approach:there is a sequence of records, each if them consist of the head and after no less than one info lihe
//Hence there is need to single out the first record; see var firstRecord bool
//It scans the "hystory.txt" and return heads of hystory records
//Errors
//1 - 1$)hyst.FindHeads:Open err=%s
//2 - 2$)hyst.FindHeads:regexp.MatchString err=%s
//3 - 3$)hyst.FindHeads:Setting currHead to zero err=%s
//4 - 4$)hyst.FindHeads:No info lines in prevHead
//5 - 5$)hyst.FindHeads:There is errr while scanning;err=%s;lineCount=%d
//6 - 6$)hyst.FindHeads:No info lines in currHead(a case of only record), line%d
//7 - 7$)hyst.FindHeads: first line must be matched; mayby there is bad pattern?, hystpattern=%v;
func FindHeads() (heads []*Head, err error) {
	var file *os.File
	var line string
	var lineCount int
	var matched bool
	var prevHead, currHead *Head //The previous and current Headers
	var onlyFirstRecord bool = true

	fmt.Printf("hyst.FindHeads:very begin	\n----end debug\n")

	file, err = os.Open("hystory.txt")
	if err != nil {
		err = fmt.Errorf("1$)hyst.FindHeads:Open err=%s", err.Error())
		return
	}
	defer file.Close()

	prevHead = new(Head)
	currHead = new(Head)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//if lineCount == 1 {
		//	fmt.Printf("hyst.FindHeads:after for scanner.Scan();in very begin; lineCount=%v\nline=%v\nmathed=%v\n----end debug\n", lineCount, scanner.Text(), matched)
		//}
		line = scanner.Text()
		lineCount++
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == byte('#') {
			continue
		}
		//firstly an attempt to try on the pattern
		if matched, err = regexp.MatchString(hystpattern, line); err != nil {
			err = fmt.Errorf("2$)hyst.FindHeads:regexp.MatchString err=%s", err.Error())
			return
		}
		if matched { //A line matches and there is not error of matching
			fmt.Printf("hyst.FindHeads: linecount=%v\nline=%v\nhystpattern=%v\n----end debug\n", lineCount, line, hystpattern)
			if err = currHead.SetZero(); err != nil {
				err = fmt.Errorf("3$)hyst.FindHeads:Setting currHead to zero err=%s", err.Error())
				return
			}
			currHead.H = line
			currHead.Lbeg = lineCount
			fmt.Printf("FindHeads:currHead:%v\n----\nonlyFirstRecord=%v\n", currHead, onlyFirstRecord)
			prevHead = currHead
			if prevHead.Lend == 0 {
				err = fmt.Errorf("4$)hyst.FindHeads:No info lines in prevHead, line%d", lineCount)
				return
			}
			if lineCount > 1 {
				onlyFirstRecord = false
			}
		} else { //not matched - first line must be matced!
			if lineCount == 1 {
				err = fmt.Errorf("hyst.FindHeads: first line must be matched; mayby there is bad pattern?, hystpattern=%v;", hystpattern)
			}
		}

		//Further we have deel with a info line of the currNead
		currHead.Lend = currHead.Lend + 1
	} //for scanner.Scan()

	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("5$)hyst.FindHeads:There is errr while scanning;err=%s;lineCount=%d", err.Error(), lineCount)
		return
	}
	if onlyFirstRecord {
		if currHead.Lend == 0 {
			err = fmt.Errorf("6$)hyst.FindHeads:No info lines in currHead(a case of only record), line%d", lineCount)
			return
		}
		heads = append(heads, currHead)
	}

	for _, v := range heads {
		fmt.Printf("hyst.FindHeads:%v\n", v.H)
		fmt.Printf("hyst.FindHeads:%v\n", v.Lbeg)
		fmt.Printf("hyst.FindHeads:%v\n", v.Lend)
		fmt.Printf("hyst.FindHeads:----end debug\n")
	}
	return
}

//Пн авг 14 11:15:45 MSK 2023
//Errors:
//1 - 1$)hyst.LastHystRec: Open file err=%s
//2 - 2$)hyst.LastHystRec: err=%s
func LastHystRec() (rec []string, err error) {
	var file *os.File
	var line string
	var lineCount int
	var heads []*Head
	var lastHead *Head

	file, err = os.Open("hystory.txt")
	if err != nil {
		err = fmt.Errorf("1$)hyst.LastHystRec: Open file err=%s", err.Error())
		return
	}
	defer file.Close()

	if heads, err = FindHeads(); err != nil {
		err = fmt.Errorf("2$)hyst.LastHystRec: err=%s", err.Error())
		return
	}
	lastHead = heads[len(heads)-1]
	//fmt.Printf("hyst.LastHystRec:lastHead:%v\n----\n", lastHead)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		lineCount++
		if lineCount <= lastHead.Lend && lineCount >= lastHead.Lbeg {
			rec = append(rec, line)
		}
	}
	return
}

func TestLastHystRec(t *testing.T) {
	lastRec, err := LastHystRec()
	if err != nil {
		fmt.Println("LastHystRec err=", err.Error())
	}
	fmt.Println("lastRec:", lastRec)
}
