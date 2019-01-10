//This is the R-available wrapper for the NY phil package
package main

// #define USE_RINTERNALS
// #include "../Rgo/sexp/Rheader/Rinternals.h"
import "C"

import(
	"github.com/EMurray16/Rgo/sexp"
	"unsafe"
	"github.com/EMurray16/coperformance/nyphil"
	//used to reformat composer names
	"strings"
	//used for troubleshooting
	//"fmt"
)

//this function is just a dereference wrapper for GoSEXP objects
func derefSEXP(s sexp.GoSEXP) C.SEXP {
	return *(*C.SEXP)(s.Point)
}

//export GetCoperformance
func GetCoperformance(min, years, fname, jsonname C.SEXP) (C.SEXP) {
	//This gets the coperformance matrix and returns the names and matrix to R
	
	//create a dummy matrix output
	f1 := make([][]float64, 3)
	for ind,_ := range f1 {
		f1[ind] = make([]float64, 3)
	}
	//create the three SEXP objects that will eventually feed into the list output
	matSEXP := sexp.Float2sexp(sexp.Matrix{f1}.Vectorize())
	namesSEXP := sexp.String2sexp("none")
	errsSEXP := sexp.String2sexp("nil")
	
	//create the sexp package aliases for the inputs
	minPoint := sexp.GoSEXP{unsafe.Pointer(&min)}
	yPoint := sexp.GoSEXP{unsafe.Pointer(&years)}
	filePoint := sexp.GoSEXP{unsafe.Pointer(&fname)}
	jsonPoint := sexp.GoSEXP{unsafe.Pointer(&jsonname)}
	
	//extract the SEXP inputs into useful values
	minProg := sexp.AsInts(minPoint)[0] //make this an int instead of a []int of length 1
	ystring := sexp.AsString(yPoint)
	infile := sexp.AsString(filePoint)
	jfile := sexp.AsString(jsonPoint)
	
	//start by getting the map of composers
	CompMap, n, err := nyphil.LoadComposers(infile, minProg, ystring)
	if err != nil {
		errsSEXP = sexp.String2sexp(err.Error())
		//now format into a list
		list := sexp.List{[]sexp.GoSEXP{matSEXP,namesSEXP,errsSEXP}}
		//create the sexp object (aka R list)
		Rlist := sexp.List2sexp(list)
		//return the derefed Rlist
		return derefSEXP(Rlist)
	}
	
	//now create the coperformance matrix
	CoperfMat, CompNames := nyphil.Coperformance(CompMap, n)
	
	//loop through the names and change all instances of 2 spaces to 1
	for ind,cname := range CompNames {
		CompNames[ind] = strings.Replace(cname, "  ", " ", -1)
	}
	
	//now format everything to be SEXPable
	CoperfMat2 := sexp.Matrix{CoperfMat}.Vectorize()
	CompNames2 := sexp.Slice2string(CompNames, "~")
	
	//write the json data file
	err = nyphil.WriteComposer(CompMap, jfile)
	if err != nil {
		errsSEXP = sexp.String2sexp("WARNING: " + err.Error())
	}
	
	//convert the data itself to sexp
	matSEXP = sexp.Float2sexp(CoperfMat2)
	namesSEXP = sexp.String2sexp(CompNames2)
	
	//now format into a sexp.List and then an R list
	list := sexp.List{[]sexp.GoSEXP{matSEXP,namesSEXP,errsSEXP}}
	Rlist := sexp.List2sexp(list)
	
	//return the derefed Rlist
	return derefSEXP(Rlist)
}

//export ProcessPhil
func ProcessPhil(min, years, inName, jsonName, csvName C.SEXP) (C.SEXP) {
	//create sexp package aliases for all the inputs
	minPoint := sexp.GoSEXP{unsafe.Pointer(&min)}
	yPoint := sexp.GoSEXP{unsafe.Pointer(&years)}
	inPoint := sexp.GoSEXP{unsafe.Pointer(&inName)}
	jsonPoint := sexp.GoSEXP{unsafe.Pointer(&jsonName)}
	csvPoint := sexp.GoSEXP{unsafe.Pointer(&csvName)}
	
	//extract the integer minimum number of programs
	minProg := sexp.AsInts(minPoint)
	//extract the filenames strings
	ystring := sexp.AsString(yPoint)
	infile := sexp.AsString(inPoint)
	jsonfile := sexp.AsString(jsonPoint)
	csvfile := sexp.AsString(csvPoint)
	
	//Now run the nyphil package workhorse
	jm, cm, lm := nyphil.ProcessJSON(minProg[0], ystring, infile, jsonfile, csvfile)
	
	//the errors need to be reformatted into a string
	ErrString := jm + " /AND/ " + cm + " /AND/ " + lm
	
	//now coerce the errors and messages to GoSEXP objects
	jpoint := sexp.String2sexp(ErrString)
	
	//dereference the unsafe pointers to normal C pointers
	Emessage := derefSEXP(jpoint)
	
	//now we can return
	return Emessage
}

func main() {}
	