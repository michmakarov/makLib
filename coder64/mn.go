//Чт мая  4 17:55:09 MSK 2023
package coder64

import (
	"encoding/base64"
	//"base58"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

const maxFileVol = 100000000 //100Mb

var eCode uint = 1          //effects the func SysErrr; 1 - default value; others state by func SetExitCode
var printStack bool         //effects the func SysErrr; fals - default value; others state by func SetTraceFlag
var print_gebug bool = true //effects the func printDebug

func QQQ() {
	fmt.Printf("The coder64.QQQ is here!\n")
}

//Чт мая  4 17:55:09 MSK 2023
//For description see func CodeFile_help
func CodeFile(origName string) {
	var err error
	var origF, resF *os.File
	var origCont []byte
	var readCount int
	var writeCount int
	var codedCont []byte
	var codedComtStr string
	var origContMimeType string

	var enc *base64.Encoding

	if origF, err = os.Open(origName); err != nil {
		SysErr(fmt.Sprintf("coder64.CodeFile: open err %s", err.Error()))
	}
	origCont = make([]byte, maxFileVol+1)
	if readCount, err = origF.Read(origCont); err != nil {
		SysErr(fmt.Sprintf("coder64.CodeFile: reading content err %s", err.Error()))
	}
	if readCount > maxFileVol {
		SysErr(fmt.Sprintf("coder64.CodeFile: file more than %d bytes in volume not allowed", maxFileVol))
	}
	origCont = origCont[0:readCount]

	//Now we have the origCont

	enc = base64.StdEncoding
	enc = enc.WithPadding('=')
	//enc = base64.NewEncoding()

	codedCont = make([]byte, base64.StdEncoding.EncodedLen(len(origCont)))
	//base64.StdEncoding.Encode(codedCont, origCont)// for edification
	enc.Encode(codedCont, origCont)
	PrintDebug(fmt.Sprintf("CodeFile:codedCont %v...%v", string(codedCont[:50]), string(codedCont[len(codedCont)-50:len(codedCont)])))
	PrintDebug(fmt.Sprintf("CodeFile:lengs %v--%v %v", len(origCont), len(codedCont), float64(len(codedCont))/float64(len(origCont))))

	origContMimeType = http.DetectContentType(origCont)
	if origContMimeType == "application/octet-stream" {
		SysErr(fmt.Sprintf("coder64.CodeFile: ", maxFileVol))
	}

	PrintDebug(fmt.Sprintf("CodeFile:origContMimeType %v", origContMimeType))
	codedComtStr = "<img alt= \"qqq\" src=\"data:" + origContMimeType + ";base64," + string(codedCont) + "\"/>"
	if resF, err = os.Create(origName + ".html"); err != nil {
		SysErr(fmt.Sprintf("coder64.CodeFile: create err %s", err.Error()))
	}
	if writeCount, err = resF.WriteString(codedComtStr); err != nil {
		SysErr(fmt.Sprintf("coder64.CodeFile: writing encoded content err %s", err.Error()))
	}
	if writeCount != len(codedComtStr) {
		SysErr(fmt.Sprintf("coder64.CodeFile: writeCount (%d) != len(codedComtStr) (%d)", writeCount, len(codedComtStr)))
	}
}

//It prints the msg and halt the application
func SysErr(msg string) {
	fmt.Printf("System error:%s\n", msg)
	if printStack {
		debug.PrintStack()
	}
	os.Exit(int(eCode))
}
func SetExitCode(code uint) {
	eCode = code
}

func SetTraceFlag(f bool) {
	printStack = f
}

func PrintDebug(msg string, forcePrint ...bool) {
	var old_print_gebug bool
	if len(forcePrint) > 0 {
		old_print_gebug = print_gebug
		print_gebug = true
		defer SetPrintDebug(old_print_gebug)
	}
	if !print_gebug {
		return
	}
	fmt.Printf("!!!Lib_DEBAG:%s\n", msg)
}

func SetPrintDebug(new bool) {
	print_gebug = new
}
