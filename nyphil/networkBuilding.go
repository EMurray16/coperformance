package nyphil

import (
	"github.com/EMurray16/Rgo/rsexp"
)

func MakeNodesAndEdges(Coperf rsexp.Matrix, names []string) (from, to []string, strengths []float64, composerCounts []int) {
	// we want to count the elements of the matrix > 1 so we know how big we can pre-allocate the output slices
	var length int = -Coperf.Nrow // this will offset the diags, which we don't want to count
	for _, f := range Coperf.Data {
		if f > 1 {
			length++
		}
	}

	from = make([]string, length)
	to = make([]string, length)
	strengths = make([]float64, length)
	composerCounts = make([]int, Coperf.Nrow)

	var ind int = 0
	for row := 0; row < Coperf.Nrow; row++ {
		for col := 0; col <= row; col++ {
			strength, _ := Coperf.GetInd(row, col)
			if row == col {
				composerCounts[row] = int(strength)
				continue
			}

			if strength > 1 {
				from[ind] = names[row]
				to[ind] = names[col]
				strengths[ind] = strength
				ind++
			}

		}
	}

	return from, to, strengths, composerCounts
}
