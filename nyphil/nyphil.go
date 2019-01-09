//This contains functions to process the JSON files for the NY philharmonic data
package nyphil

import(
	//"fmt"
	//used to read/write json
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
	//used for the coperformance matrix
	"github.com/EMurray16/Rgo/rfunc"
	"math"
	//used for writing/reading csv files
	"encoding/csv"
	"io"
	"strings"
	//used for error handling
	"errors"
)

var (
	YearInds = loadYears()
	NoDataErr = errors.New("Not enough data to write file")
)

//this function loads a map of seasons to indexes
func loadYears() (map[string]int) {
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

//create a composer type
type Composer struct {
	Nprog int
	Programs []uint16
	Seasons []string
}
	
//this function returns a map of composer names to program ids and seasons
func LoadComposers(file string, MinProg int, Years string) (map[string]Composer, int, error) {
	//convert the years string to a min and max year
	seasonbounds := strings.Split(Years, ":")
	boundearly := YearInds[ seasonbounds[0] ]
	boundlate := YearInds[ seasonbounds[1] ]
	
	//start by making the map, guessing they'll be about 3000 composers
	Pmap := make(map[string]([]uint16), 3000)
	Smap := make(map[string]([]string), 3000)
	OutMap := make(map[string]Composer, 300)
	
	//now open the json file and read it into an interface
	var i interface{}
	
	f, err := os.Open(file)
	if err != nil {
		return OutMap, 0, err
	}
	asbytes, err := ioutil.ReadAll(f)
	if err != nil {
		return OutMap, 0, err
	}
	err = json.Unmarshal(asbytes, &i)
	if err != nil {
		return OutMap, 0, err
	}
	
	//pull the slice of programs out of the interface
	programs := i.(map[string]interface{})["programs"].([]interface{})
	N := 0
	
	//loop through the programs to pull out the program id
	for _,program := range programs {
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
			return OutMap, N, err
		}

		//now isolate the works
		works := progmap["works"].([]interface{})
		//loop through the works and get a slice of composer names
		nameslice := make([]string, 0)
		for _,work := range works {
			nameface := work.(map[string]interface{})["composerName"]
			if nameface != nil {
				nameslice = append(nameslice, nameface.(string))
			}
		}
		
		//now get the unique composer names
		UniqName := rfunc.StringUnique(nameslice)
		
		//add the program to the slice for each name
		for _,name := range UniqName {
			if Pmap[name] == nil {
				Pmap[name] = []uint16{}
				//the seasons map will also be nil
				Smap[name] = []string{}
			}
			Pmap[name] = append(Pmap[name], uint16(progid))
			Smap[name] = append(Smap[name], season)
		}
	}
	
	//now that we have all composers, only keep those with minimum program size
	for name,progs := range Pmap {
		if len(progs) >= MinProg {
			OutMap[name] = Composer{len(progs), progs, Smap[name]}
		}
	}
	
	//when we get here, we're done and can return
	return OutMap, N, nil
}

//this function makes the coperformance matrix from a map of composers
func Coperformance(CompMap map[string]Composer, Nprog int) ([]([]float64), []string) {
	//we'll make a dummy slice of composers to range faster and more flexibly
	Ncomp := len(CompMap)
	if Ncomp == 0 {
		return []([]float64){}, []string{"0"}
	}
	
	CompSlice := make([]Composer, Ncomp)
	CompNames := make([]string, Ncomp)
	i := 0
	for name,comp := range CompMap {
		CompSlice[i] = comp
		CompNames[i] = name
		i++
	}
	
	//now we need to build coperformance matrix
	CoperfMat := make([]([]float64), Ncomp)
	//initialize the 0th column and the first diag element
	CoperfMat[0] = make([]float64, Ncomp)
	//now loop through the composers to index the rows
	for row := 1; row < len(CompMap); row++ {
		//initialize the column of the same index
		CoperfMat[row] = make([]float64, Ncomp)
		
		//now range through all the composers up to the diag
		for col := 0; col < (row+1); col++ {
			//if the col == row, the diag is defined as the number of appearances
			if col == row {
				CoperfMat[row][col] = float64(CompSlice[row].Nprog)
				continue
			}
			//isolate the two composers
			RowComp := CompSlice[row]
			ColComp := CompSlice[col]
			
			//find the denominator of coperformance
			denominator := math.Ceil(float64(RowComp.Nprog * ColComp.Nprog) / float64(Nprog))
			//the numerator is the number of programs in common
			numerator := float64(Ncommon_Uint16(RowComp.Programs, ColComp.Programs))
			
			//put the element in the matrix
			CoperfMat[row][col] = numerator / denominator
		}
	}
	//we need to add the frist element of the diag now that its initialized
	CoperfMat[0][0] = float64(CompSlice[0].Nprog)
	
	return CoperfMat, CompNames
}

//this function writes the json of Composers given a map[string]Composers
func WriteComposer(CompMap map[string]Composer, filename string) (error) {
	//write the composer info out to a json
	jsonbytes, err := json.Marshal(CompMap)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, jsonbytes, os.ModeDir)
	if err != nil {
		return err
	}
	
	return nil
}

//this function writes the coperformance matrix
func WriteCoperformance(Coperf [][]float64, CompNames []string, filename string) (error) {
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
	for i,name := range CompNames {
		//convert the slice of column values to strings
		Costrings := make([]string, len(CompNames))
		for j := 0; j < len(CompNames); j++ {
			Costrings[j] = strconv.FormatFloat(Coperf[i][j], 'f', -1, 64)
		}
		//now make the full row slice
		rowslice := append([]string{name}, Costrings...)
		matwriter.Write(rowslice)
	}
	
	return nil
}

//now export a function that does everything
func ProcessJSON(minprog int, yearbounds, infile, jsonfile, csvfile string) (jsonMessage, csvMessage, loadMessage string) {
	//set the error strings to "nil"
	jsonMessage = "nil"
	csvMessage = "nil"
	loadMessage = "nil"
	
	//call loadComposers
	CompMap, n, err := LoadComposers(infile, minprog, yearbounds)
	if err != nil {
		loadMessage = err.Error()
		//this is catastrophic so return
		return jsonMessage, csvMessage, loadMessage
	} 
	
	//now write the json file
	jsonerr := WriteComposer(CompMap, jsonfile)
	if jsonerr != nil {
		jsonMessage = jsonerr.Error()
	}
	
	//now make the coperformance matrix
	CoMat, names := Coperformance(CompMap, n)
	
	//now write the coperformance matrix csv
	csverr := WriteCoperformance(CoMat, names, csvfile)
	if csverr != nil {
		csvMessage = csverr.Error()
	}
	
	//now return everything
	return jsonMessage, csvMessage, loadMessage
}
	
	
