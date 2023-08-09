//Пн июл 31 16:02:13 MSK 2023
package maklib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type setVar func(optName string, optValue string) (err error)

var ReadOptDescr = `
The func ReadOpt(f setVar)(setted []string, err error) take as its argument the func of type func(optName string, optValue string) (err error)
It look up in working directory for a text file the with name of optons.txt. and scans it.
If the func find a sifnificant line it extract an option name amd value from the line
and invoke its parameter with them and if it retuned nil error adds the string of "n+"<-"+v" to the setted.
Here the n is a number of a line in the optons.txt (first line is one)
and a v is an option value
FORMAT of aa sifnificant line:
<Option name> = <OptionValue>
or
<Option name> = <OptionValue>#<comment>
Where:
<Option name> and <OptionValue> are are strings without the blanks, =, #
<comment> - on arbitrary string
ABOUT function setVar.
Of cause, on a calling side it maybe any function that takes two strings and returrn an error.
Only it must be strongly taken in account that if no a significant lines the func is not be colled.
The developer suppose that the above info gives enough the additional info for creating the appropriate algorithm.
`

func ReadOpt(f setVar) (setted []string, err error) {
	//var err error
	var file *os.File
	var line string
	var opt []string //the option name id otp[0]; the optiom value (maybe with a comment) ig opt[1]
	var optionName, optionValue string
	var optWithlineComment []string
	var lineCount int

	//var setted []string
	var addToSetted = func(n, v string) {
		setted = append(setted, n+"<-"+v)
	}

	file, err = os.Open("options.txt")
	if err != nil {
		err = fmt.Errorf("readOp", err.Error())
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
		opt = strings.Split(line, "=")
		if len(opt) != 2 {
			err = fmt.Errorf("There is bad format of a line : there is not \"=\" or more than one such\n %s", line)
			err = fmt.Errorf("readOp", err.Error())
			return
		} else {
			optWithlineComment = strings.Split(opt[1], "#")
			if len(optWithlineComment) == 1 { //no comment
				optionName = strings.TrimSpace(opt[0])
				optionValue = strings.TrimSpace(opt[1])
			} else {
				optionName = strings.TrimSpace(opt[0])
				optionValue = strings.TrimSpace(optWithlineComment[0])
			}
		}

		//
		if err = f(optionName, optionValue); err != nil {
			err = fmt.Errorf("readOp: setVar err=", err.Error())
			return
		} else {
			addToSetted(fmt.Sprint(lineCount), optionName)
		}

		if err = scanner.Err(); err != nil {
			err = fmt.Errorf("There is errr while scanning the option.txt; lineCount=%d", lineCount)
			return
		}
	} //for scanner.Scan()
	return
}
