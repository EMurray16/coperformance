package nyphil

//this function returns the number of common occurrences in 2 slices of uints
func Ncommon_Uint16(Slice1, Slice2 []uint16) (common int) {
	//compare all elements of both slices, add when common
	for _,i1 := range Slice1 {
		for _,i2 := range Slice2 {
			if i1 == i2 {
				common++
			}
		}
	}
	return common
}