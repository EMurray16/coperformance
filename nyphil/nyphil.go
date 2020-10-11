//This contains functions to process the JSON files for the NY philharmonic data
package nyphil

import (
	//"fmt"
	//used to read/write json
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	//used for the coperformance matrix
	"github.com/EMurray16/Rgo/rfunc"
	"github.com/EMurray16/Rgo/rsexp"
	"math"
	//used for writing/reading csv files
	"encoding/csv"
	"io"
	"strings"
)

//this function loads a map of seasons to indexes
func loadYears() map[string]int {
	OutMap := make(map[string]int, 175)
	//open the file and make the reader
	f, _ := os.Open("./philyears.csv")
	r := csv.NewReader(f)

	for i := 0; ; i++ {
		year, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return OutMap
		}

		OutMap[year[0]] = i
	}

	return OutMap
}

// loadComposers creates a map of composer names to their programs and the number of total programs
func LoadComposers(file string, MinProg int, Years string) (composers map[string]composer, nPrograms int, err error) {
	//convert the years string to a min and max year
	seasonbounds := strings.Split(Years, ":")
	boundearly := YearInds[seasonbounds[0]]
	boundlate := YearInds[seasonbounds[1]]

	//start by making the map, guessing they'll be about 3000 composers
	programMap := make(map[string]([]uint16), 3000)
	seasonMap := make(map[string]([]string), 3000)
	composers = make(map[string]composer, 3000)

	//now open the json file and read it into an interface
	var i interface{}

	f, err := os.Open(file)
	if err != nil {
		return composers, 0, err
	}
	asbytes, err := ioutil.ReadAll(f)
	if err != nil {
		return composers, 0, err
	}
	err = json.Unmarshal(asbytes, &i)
	if err != nil {
		return composers, 0, err
	}

	//pull the slice of programs out of the interface
	programs := i.(map[string]interface{})["programs"].([]interface{})
	N := 0

	//loop through the programs to pull out the program id
	for _, program := range programs {
		progmap := program.(map[string]interface{})

		//check the season
		season := progmap["season"].(string)
		if YearInds[season] < boundearly || YearInds[season] > boundlate {
			continue
		}
		//index the N counter
		N++

		//get the id, which is a string, and parse to int
		id_string := progmap["programID"].(string)
		progid, err := strconv.ParseUint(id_string, 10, 16)
		if err != nil {
			return composers, N, err
		}

		//now isolate the works
		works := progmap["works"].([]interface{})
		//loop through the works and get a slice of composer names
		nameslice := make([]string, 0)
		for _, work := range works {
			nameface := work.(map[string]interface{})["composerName"]
			if nameface != nil {
				nameslice = append(nameslice, nameface.(string))
			}
		}

		//now get the unique composer names
		UniqName := rfunc.StringUnique(nameslice)

		//add the program to the slice for each name
		for _, name := range UniqName {
			if programMap[name] == nil {
				programMap[name] = []uint16{}
				//the seasons map will also be nil
				seasonMap[name] = []string{}
			}
			programMap[name] = append(programMap[name], uint16(progid))
			seasonMap[name] = append(seasonMap[name], season)
		}
	}

	//now that we have all composers, only keep those with minimum program size
	for name, progs := range programMap {
		if len(progs) >= MinProg {
			composers[name] = composer{len(progs), progs, seasonMap[name]}
		}
	}

	//when we get here, we're done and can return
	return composers, N, nil
}

//this function makes the coperformance matrix from a map of composers
func Coperformance(CompMap map[string]composer, Nprog int) (rsexp.Matrix, []string, error) {
	//we'll make a dummy slice of composers to range faster and more flexibly
	Ncomp := len(CompMap)
	if Ncomp == 0 {
		return rsexp.Matrix{}, nil, NoDataErr
	}

	CompSlice := make([]composer, Ncomp)
	CompNames := make([]string, Ncomp)
	i := 0
	for name, comp := range CompMap {
		CompSlice[i] = comp
		CompNames[i] = name
		i++
	}

	//now we need to build coperformance matrix
	CoperfMat, err := rsexp.CreateIdentity(Ncomp)
	if err != nil {
		return rsexp.Matrix{}, CompNames, err
	}

	//now loop through the composers to index the rows
	for row := 1; row < len(CompMap); row++ {

		//now range through all the composers up to the diag
		for col := 0; col < (row + 1); col++ {
			//if the col == row, the diag is defined as the number of appearances
			if col == row {
				err := CoperfMat.SetInd(row, col, float64(CompSlice[row].Nprog))
				if err != nil {
					return *CoperfMat, CompNames, err
				}
				continue
			}
			//isolate the two composers
			RowComp := CompSlice[row]
			ColComp := CompSlice[col]

			//find the denominator of coperformance
			denominator := math.Ceil(float64(RowComp.Nprog*ColComp.Nprog) / float64(Nprog))
			//the numerator is the number of programs in common
			numerator := float64(Ncommon_Uint16(RowComp.Programs, ColComp.Programs))

			//put the element in the matrix
			err := CoperfMat.SetInd(row, col, numerator/denominator)
			if err != nil {
				return *CoperfMat, CompNames, err
			}
		}
	}
	//we need to add the frist element of the diag now that its initialized
	err = CoperfMat.SetInd(0, 0, float64(CompSlice[0].Nprog))
	if err != nil {
		return *CoperfMat, CompNames, err
	}

	return *CoperfMat, CompNames, nil
}

//this function writes the json of Composers given a map[string]Composers
func WriteComposers(CompMap map[string]composer, filename string) error {
	jsonbytes, err := json.Marshal(CompMap)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(jsonbytes)
	if err != nil {
		return err
	}

	return nil
}

//this function writes the coperformance matrix
func writeCoperformance(Coperf rsexp.Matrix, CompNames []string, filename string) error {
	if len(CompNames) < 2 {
		return NoDataErr
	}
	//now write a csv of the Coperformance matrix
	mfile, err := os.Create(filename)
	defer mfile.Close()
	if err != nil {
		return err
	}

	//make a writer and write the first row
	matwriter := csv.NewWriter(mfile)
	defer matwriter.Flush()
	matwriter.Write(append([]string{""}, CompNames...))

	//now write each row of the file
	for i, name := range CompNames {
		//convert the slice of column values to strings
		rowslice := make([]string, len(CompNames)+1)
		rowslice[0] = name
		coperfRow, err := Coperf.GetRow(i)
		if err != nil {
			return err
		}

		for j := 0; j < len(CompNames); j++ {
			rowslice[j+1] = strconv.FormatFloat(coperfRow[j], 'f', -1, 64)
		}

		//now make the full row slice
		matwriter.Write(rowslice)
	}

	return nil
}