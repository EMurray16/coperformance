package main

// #include <Rinternals.h>
// #cgo CFLAGS: -I/Library/Frameworks/R.framework/Headers/
// #cgo CFLAGS: -I/usr/share/R/include/
import "C"

import (
	"coperformance/nyphil"
	"errors"
	"github.com/EMurray16/Rgo/rsexp"
)

//this function is just a dereference wrapper for GoSEXP objects
func deref(s rsexp.GoSEXP) C.SEXP {
	return *(*C.SEXP)(s.Point)
}

// we will create package level variables for the outputs so they're usable in all functions
var (
	coperfMat  = rsexp.Matrix{}
	composerNames = make([]string, 0)
	networkList = rsexp.List{}
	nilError = errors.New("nil")
)

func formatOutput(errData, errFile error) C.SEXP {
	errString := errData.Error()
	errString2 := errFile.Error()

	outList := rsexp.NewList(rsexp.Matrix2sexp(coperfMat),
		rsexp.String2sexp(composerNames),
		rsexp.List2sexp(networkList),
		rsexp.String2sexp([]string{errString}),
		rsexp.String2sexp([]string{errString2}))

	return deref(rsexp.List2sexp(outList))
}

//export ProcessNYPhilJSON
func ProcessNYPhilJSON(filenames, minPerf, years C.SEXP) C.SEXP {
	// we return a list with the following objects:
		// 1. The coperformance matrix
		// 2. A vector of composer names
		// 3. A list with the info needed to make the network
		// 4. A string with an error message if there is one
	// parse the input into Go variables
	filenameGS, err := rsexp.NewGoSEXP(&filenames)
	if err != nil {
		return formatOutput(err, err)
	}
	// filenameGS := rsexp.GoSEXP{unsafe.Pointer(&filename)}
	filenameString, err := filenameGS.AsStrings()
	if err != nil {
		return formatOutput(err, err)
	}

	minPerfGS, err := rsexp.NewGoSEXP(&minPerf)
	if err != nil {
		return formatOutput(err, err)
	}
	minPerfInt, err := minPerfGS.AsInts()
	if err != nil {
		return formatOutput(err, err)
	}

	yearsGS, err := rsexp.NewGoSEXP(&years)
	if err != nil {
		return formatOutput(err, err)
	}
	yearsString, err := yearsGS.AsStrings()
	if err != nil {
		return formatOutput(err, err)
	}

	// start by loading the composers
	composerMap, nProg, err := nyphil.LoadComposers(filenameString[0], minPerfInt[0], yearsString[0])
	if err != nil {
		return formatOutput(err, err)
	}

	// now build the coperformance matrix
	coperfMat, composerNames, err = nyphil.Coperformance(composerMap, nProg)
	if err != nil {
		return formatOutput(err, err)
	}

	// now build the network outputs
	from, to, strength, count := nyphil.MakeNodesAndEdges(coperfMat, composerNames)
	networkList = rsexp.NewList(rsexp.String2sexp(from),
		rsexp.String2sexp(to),
		rsexp.Float2sexp(strength),
		rsexp.Int2sexp(count))

	// write the composer output to json
	err = nyphil.WriteComposers(composerMap, filenameString[1])
	if err != nil {
		return formatOutput(nilError, err)
	}

	// if we get this far, return with a nil error
	return formatOutput(nilError, nilError)
}


func main() {}
