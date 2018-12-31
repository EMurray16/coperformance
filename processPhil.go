//This is the R-available wrapper for the NY phil package
package main

// #define USE_RINTERNALS
// #include <Rinternals.h>
import "C"

/* testing code
	dyn.load("C:/Users/Evan/go/src/NYphil/processPhil.so")
	infile = as.integer(charToRaw("C:/Users/Evan/go/src/NYphil/complete.json"))
	jfile = as.integer(charToRaw("C:/evan/documents/R/testcomp.json"))
	cfile = as.integer(charToRaw("C:/evan/documents/R/testcoperf.csv"))
	minint = as.integer(500)
	years = as.integer(charToRaw("1856-57:2017-18"))
	ResList = .Call("ProcessPhil", minint, years, infile, jfile, cfile)
*/

import(
	"Rgo/sexp"
	"unsafe"
	"Nyphil/nyphil"
	//used for troubleshooting
	"io/ioutil"
	"os"
	"fmt"
)

//this function is just a dereference wrapper for GoSEXP objects
func derefSEXP(s sexp.GoSEXP) C.SEXP {
	return *(*C.SEXP)(s.Point)
}

//export ProcessPhil
func ProcessPhil(min, years, inName, jsonName, csvName C.SEXP) (C.SEXP) {
	//create sexp package aliases for all the inputs
	minPoint := sexp.GoSEXP{unsafe.Pointer(&min)}
	yPoint := sexp.GoSEXP{unsafe.Pointer(&years)}
	inPoint := sexp.GoSEXP{unsafe.Pointer(&inName)}
	jsonPoint := sexp.GoSEXP{unsafe.Pointer(&jsonName)}
	csvPoint := sexp.GoSEXP{unsafe.Pointer(&csvName)}
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/progress.txt", []byte("Converted inputs to GoSEXPs"), os.ModeDir)
	
	//extract the integer minimum number of programs
	minProg := sexp.AsInts(minPoint)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/minProg.txt", []byte(fmt.Sprint(minProg)), os.ModeDir)
	//extract the filenames strings
	ystring := sexp.AsString(yPoint)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/ystring.txt", []byte(ystring), os.ModeDir)
	infile := sexp.AsString(inPoint)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/infile.txt", []byte(infile), os.ModeDir)
	jsonfile := sexp.AsString(jsonPoint)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/jsonfile.txt", []byte(jsonfile), os.ModeDir)
	csvfile := sexp.AsString(csvPoint)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/csvfile.txt", []byte(csvfile), os.ModeDir)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/progress.txt", []byte("Converted SEXP to Go objects"), os.ModeDir)
	
	//Now run the nyphil package workhorse
	jm, cm, lm := nyphil.ProcessJSON(minProg[0], ystring, infile, jsonfile, csvfile)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/progress.txt", []byte("Ran nyphil workhorse"), os.ModeDir)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/jm.txt", []byte(jm), os.ModeDir)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/cm.txt", []byte(cm), os.ModeDir)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/em.txt", []byte(lm), os.ModeDir)
	
	//the errors need to be reformatted into a string
	ErrString := fmt.Sprint(jm, " /AND/ ", cm, " /AND/ ", lm)
	
	//now coerce the errors and messages to GoSEXP objects
	jpoint := sexp.String2sexp(ErrString)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/progress.txt", []byte("Made new GoSEXPs"), os.ModeDir)
	
	//dereference the unsafe pointers to normal C pointers
	Emessage := derefSEXP(jpoint)
	ioutil.WriteFile("C:/Users/Evan/go/src/NYphil/progress.txt", []byte("Dereferenced new GoSEXPs"), os.ModeDir)
	
	//now we can return
	return Emessage
}

func main() {}
	