// Чт мая  4 19:45:04 MSK 2023
package coder64

//import (
//	"fmt"
//)

func CodeFile_help() string {
	return `
func CodeFile(origName string)
The codeFile function produces HTML tag <image src="data: ..> as a text tag file
where the data is the 64 encode of the content of file with name passing as paramrter origName,
for examle LUBA.ipg .
The result text file has a name of fofmat: <origName>.html, foe example LUBA.ipg.html .
If it encounters an error it halts the application.
The function seaks for original image files the local directory "images" that must be in the working directory,
It places the result in the same spot.
`
}

func Package_help() string {
	return `
The coder64 package for io operations is using only the os package
in hope that the package allows concurrent working in many goroutenes.
At Пт мая  5 14:21:53 MSK 2023 the author does not know if his hope is enough justified.
So he is going to write a test for heavy testing.
`
}
