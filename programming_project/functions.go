package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"math"
	"math/rand"
)
 
/*
---------------------------------
IN/OUT FUNCTIONS
---------------------------------
*/
//we need to establish a datatype set
//ok, and to do this we will work in strings. just for simplicity's sake.
type dataframe [][]string

//for ease of manipulation later, I'm goibng to be adding a numeric dataframe type too.
type dfNumeric struct {
	values [][]float64
	names []string
}

//ReadCSV takes in a string name that will be the filename for the CSV we want to read. 
//It returns a dataframe where each sub-array in the array corresponds to a row of the CSV.
func ReadCSV(filename string) dataframe {
	var df dataframe
	currentIndex := 0

	//standard file reading code
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentRow := scanner.Text()
		entries := strings.Split(currentRow, ",")

		//make a new row
		newRow := make([]string, len(entries))
		df = append(df, newRow)
		for i := range entries {
			df[currentIndex][i] = entries[i]
		}

		currentIndex++
	}

	err2 := scanner.Err()
	if err2 != nil {
		panic(err2)
	}

	return df
}

//ReadNumericCSV takes in a file name and reads in only the numeric columns to a 
//NumericDF type. Critically, it assumes that all the entries aside from the first row are numerical. 
func ReadNumericCSV(filename string) dfNumeric {
	var df dfNumeric
	currentIndex := 0

	//standard file reading code
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentRow := scanner.Text()
		entries := strings.Split(currentRow, ",")

		//first case: our column names
		if currentIndex == 0 {
			newRow := make([]string, len(entries))
			for i := range entries {
				newRow[i] = entries[i]
			}

			df.names = newRow
		} else { //read in the numeric things
			newRow := make([]float64, len(entries))
			for i := range entries {
				val, err3 := strconv.ParseFloat(entries[i], 64)
				if err3 != nil {
					panic(err3)
				}
				newRow[i] = val
			}

			df.values = append(df.values, newRow)

		}

		//make a new row
		
		
		currentIndex++
	}

	err2 := scanner.Err()
	if err2 != nil {
		panic(err2)
	}

	return df
}

//WriteCSV takes in a dataframe and a file name and writes a CSV file of that name.
func WriteCSV(df dataframe, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic("oops another error")
	}

	writer := bufio.NewWriter(file)

	for i := range df {
		for j := range df[i] {
			n := len(df[i]) - 1
			fmt.Fprint(writer, df[i][j])
			if j != n {
				fmt.Fprint(writer, ",")
			}
		}
		fmt.Fprintln(writer, "")
	}

	writer.Flush()

	file.Close()
}

//WriteCorrMatrixToFile takes in a file name, matrix, names, and a filename and 
//writes a file with the proper correlations.
//Note: This is borrowed from the Metagenomics recitation; it is not particularly helpful
//in the context of a report as we will make a fancy graphic in R, but it helps for 
//knowing the actual values (can be helpful in report.)
func WriteCorrMatrixToFile(mtx [][]float64, cols []string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic("oops another error")
	}

	writer := bufio.NewWriter(file)

	//gap at start of file
	fmt.Fprint(writer, ",")

	//write the column names
	for _, col := range cols {
		fmt.Fprint(writer, col)
		fmt.Fprint(writer, ",")
	}
	fmt.Fprintln(writer, "")

	for i := range mtx {
		fmt.Fprint(writer, cols[i])
		fmt.Fprint(writer, ",")
		for j := range mtx[i] {
			fmt.Fprint(writer, mtx[i][j])
			fmt.Fprint(writer, ",")
		}
		fmt.Fprintln(writer, "")
	}

	writer.Flush()

	file.Close()
}

/*
---------------------------------
DATA CLEANING AND MANIPULATION
---------------------------------
*/

/*
For all of these: as a note, since I generally have no idea how a dataset might look,
I am choosing to use append() for transformations. This can be memory intensive.
Luckily, I have a computer that is NOT a potato, so this will be okay.
*/

//NumCols is a method for a dataframe that returns the number of columns.
func (df dataframe) NumCols() int {
	if len(df) < 1 {
		panic("error: empty dataframe given to NumCols()")
	}
	return len(df[0])
}

//NumRows is a method for a dataframe that returns the number of rows.
func (df dataframe) NumRows() int {
	return len(df)
}

//RemoveNAs is a dataframe method that takes in an integer index and removes all rows where we have an
//NA in that index.
func (df dataframe) RemoveNAs(index int) dataframe {
	//error statements.
	if len(df) < 1 {
		panic("empty dataframe passed to remove NAs")
	} else if len(df[0]) < 1 {
		panic("empty dataframe passed to remove NAs")
	} else if len(df[0]) < index {
		panic("index too big for RemoveNAs.")
	}

	//ok now we can actually do stuff
	var newdf dataframe
	//range through and add valid rows
	for i := range df {
		if df[i][index] != "" && df[i][index] != "NA" {
			newdf = append(newdf, df[i])
		}
	}

	return newdf
}

//GetColumn is a method for a dataframe that takes in an index and returns a column corresponding to that
//column position in the dataframe. 
func (df dataframe) GetColumn(index int) []string {
	//error statements.
	if len(df) < 1 {
		panic("empty dataframe passed to get column")
	} else if len(df[0]) < 1 {
		panic("empty dataframe passed to get column")
	} else if len(df[0]) < index {
		panic("index too big for RemoveNAs.")
	}

	// ok now we can do stuff
	col := make([]string, len(df))
	for i := range df {
		col[i] = df[i][index]
	}

	return col
}

//DropColumn is a function that takes in a dataframe and an index for a column of the frame. It returns all columns but that one.
func (df dataframe) DropColumn(index int) dataframe {
	var newdf dataframe

	//range through rows
	for i := range df {
		var newRow []string

		//range through cols
		for j := range df[i] {
			if j != index {
				newRow = append(newRow, df[i][j])
			}
		}

		newdf = append(newdf, newRow)
	}

	return newdf
}

//Select is a function that is basically the opposite of DropColumn; it takes an array of names and selects only those columns.
func (df dataframe) Select(names []string) dataframe {
	var newCols dataframe

	//range through and get all columns
	for i := range names {
		if ColIn(names[i], df) {
			newCol := df.GetColumn(df.GetIndex(names[i]))
			newCols = append(newCols, newCol)
		}
	}

	newDF := JoinCols(newCols)
	return newDF
}

//AddColOfOnes is a method for a dataframe that takes in a column name and adds in a new column of ones at the end of the dataframe.
func (df dataframe) AddColOfOnes(name string) dataframe {
	var newdf dataframe

	//range through all rows and add the 1 to the end
	//this 1 is a string.
	for i := range df {
		var newRow []string

		//range through cols
		for j := range df[i] {
			newRow = append(newRow, df[i][j])
		}

		//beginning: append the name
		if i == 0 {
			newRow = append(newRow, name)
		} else {
			newRow = append(newRow, "1")
		}
		
		newdf = append(newdf, newRow)
	}

	return newdf

}

//AddColOfOnes is a method for a dataframe that takes in a column name and adds in a new column of zeroes at the end of the dataframe.
func (df dataframe) AddColOfZeroes(name string) dataframe {
	var newdf dataframe

	//range through all rows and add the 1 to the end
	//this 1 is a string.
	for i := range df {
		var newRow []string

		//range through cols
		for j := range df[i] {
			newRow = append(newRow, df[i][j])
		}

		//beginning: append the name
		if i == 0 {
			newRow = append(newRow, name)
		} else {
			newRow = append(newRow, "0")
		}
		
		newdf = append(newdf, newRow)
	}

	return newdf
}

//ColIn is a function that returns true if an input string is in the input dataframe's first row.
//Otherwise, it returns false.
func ColIn(name string, df dataframe) bool {
	//range through all col names and compare to inpuit
	for i := range df[0] {
		if df[0][i] == name {
			return true
		}
	}

	//else, we haven't found the name, so return false
	return false
}

//JoinCols is a function that takes the columns given by Select() and GetColumn() and joiuns them
//into an appropriate dataframe.
func JoinCols(cols1 dataframe) dataframe {
	var newDF dataframe
	
	//range through and append as necessary
	for i := range cols1[0] {//first column given, which indexes the rows. Does that make sense? Kindof?
		var newRow []string

		for j := range cols1 { //range through the columns now
			newRow = append(newRow, cols1[i][j])
		}

		newDF = append(newDF, newRow)
	}

	return newDF
}

//GetIndex is a function that takes in a string value of a name for a column in a dataframe and returns the integer index of that variable.
func (df dataframe) GetIndex(name string) int {
	for i := range df[0] {
		if df[0][i] == name {
			return i
		}
	}

	panic("Error: name does not exist in dataframe given to GetIndex().")
	return -1
}

//IsNumeric is a function that takes in a slice of strings and returns whether or not
//one of them is a number or not, by using the error of ParseFloat.
func IsNumeric(col []string) bool {
	for i := range col {
		_, err := strconv.ParseFloat(col[i], 64)
		if err != nil {
			return false
		}
	}

	return true
}

//MakeNumeric is function that takes in a slice of strings and returns a slice of those
//strings converted to floats.
func MakeNumeric(col []string) []float64 {
	newCol := make([]float64, len(col))

	for i := range col {
		newVal, err := strconv.ParseFloat(col[i], 64)
		if err != nil {
			panic(err)
		}
		newCol[i] = newVal
	}

	return newCol
}

//ConvertToString takes in a numeric dataframe and converts all the entries to strings for writing.
func ConvertToString(data [][]float64) dataframe {
	var newdf dataframe

	for i := range data {
		var newRow []string
		for j := range data[i] {
			currentVal := strconv.FormatFloat(data[i][j], 'f', -1, 64)
			newRow = append(newRow, currentVal)
		}

		newdf = append(newdf, newRow)
	}

	return newdf
}

//GetRandPoints is a function that takes in an input number of points on Planet Earth to draw from.
//It returns a 2D array of floats, with each array inside corresponding to a single point.
//like [[3 4] [5 6] [7 8]] except with actual latitude and longitude
func GetRandPoints(num int) dataframe {
	var points [][]string

	//loop num times
	for i := 0; i < num; i++ {
		currentPoint := GetOneRandPoint()

		//make a new point
		var newPoint []string
		newPoint = append(newPoint, "None") //corresponds to gBIF ID
		newPoint = append(newPoint, "None") //corresponds to scientific name
		for j := range currentPoint {
			newPoint = append(newPoint, strconv.FormatFloat(currentPoint[j], 'f', -1, 64)) //appends the number
		}
		newPoint = append(newPoint, "0") //corresponds to presence/absence

		points = append(points, newPoint)
	}

	return points
}

//GetOneRandPoint is a function that takes no inputs and returns a random point with random lat, long,
//depth and elevation.
func GetOneRandPoint() []float64 {
	currentPoint := make([]float64, 4)

	//get random latitude between -90 and 90
	randLat := rand.Intn(181)
	randLat -= 90 //so we get proper range
	randFloat1 := rand.Float64()
	randLatFloat := float64(randLat) * randFloat1 //so we add a decimal

	//get random longitude between -180 and 180
	randLong := rand.Intn(361)
	randLong -= 180
	randFloat2 := rand.Float64()
	randLongFloat := float64(randLong) * randFloat2

	//get random elevation, between 1 and 500
	randElev := rand.Intn(501)

	//get random depth, between 1 and 100
	randDepth := rand.Intn(101)

	currentPoint[0] = randLatFloat
	currentPoint[1] = randLongFloat
	currentPoint[2] = float64(randElev)
	currentPoint[3] = float64(randDepth)

	return currentPoint
}

//AddHeadingToRandPoints takes in the random points made by GetRandPoints() and converted to string by ConvertToString()
//It adds a header to the top as the titles for the dataframe.
func AddHeadingToRandPoints(df dataframe) dataframe {
	var newdf dataframe
	head := make([]string, 7)

	head[0] = "gBIF ID"
	head[1] = "scientificName"
	head[2] = "decimalLatitude"
	head[3] = "decimalLongitude"
	head[4] = "elevation"
	head[5] = "depth"
	head[6] = "present"

	newdf = append(newdf, head)
	newdf = append(newdf, df...)

	return newdf
}

//SortRandPoints is a function that takes in two dataframes, df1 and df2. It returns the points of df2 filtered out so that
//the latitudes of the points are not close to the points of df1, with those points added to df1.
func SortRandPoints(df1, df2 dataframe) dataframe {
	lat1 := df1.GetColumn(2)
	lat2 := df2.GetColumn(2)
	long1 := df1.GetColumn(3)
	long2 := df2.GetColumn(3)

	//make these into numerics
	lat1_num := MakeNumeric(lat1[1:])
	lat2_num := MakeNumeric(lat2[1:])
	long1_num := MakeNumeric(long1[1:])
	long2_num := MakeNumeric(long2[1:])

	//first we make the new dataframe
	var newdf dataframe
	newdf = append(newdf, df1...)

	threshold := 5.0

	for i := range df2[1:] {
		//we'll compare every point in df2 to every point in df1
		//to see if it belongs
		inProx := false //make a new variable that determines if a point in range

		for j := range df1[1:] {
			currentDist := Distance(lat1_num[j], lat2_num[i], long1_num[j], long2_num[i])
			if currentDist < threshold {
				//ok, if we're here, we have a distance less than our threshold.
				//so we set the variable to true
				inProx = true
			}
		}

		//now that we're here, we need to determine whether or not to append the point
		//and we use the inProx variable
		if inProx == false {
			newdf = append(newdf, df2[i+1])
		}
	}

	return newdf
}

//Distance is a classic function that takes in two points and returns the distance between them.
//This is adapted from class to remove the OrderedPoint object.
func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(((x2-x1) * (x2-x1)) + ((y2-y1) * (y2-y1)))
}

//Bootstrap is a method for a dataframe that takes in an integer number of terms to bootstrap from the dataset.
//It randomly selects that many items from the dataset (repeats allowed) and adds them to the original dataframe.
func (df dataframe) Bootstrap(num int) dataframe {
	var newdf dataframe
	n := len(df)

	for i := 0; i < num; i++ {
		//choose random row
		randIndex := rand.Intn(n)
		newRow := df[randIndex]

		//add to dataframe
		newdf = append(newdf, newRow)
	}

	df = append(df, newdf...)
	return df
}

//CombineDataFrames takes two dataframes and places the first on top of the other, removing the column names of the second.
//It only works for dataframes that have the same number of columns.
func CombineDataFrames(df1, df2 dataframe) dataframe {
	if df1.NumCols() != df2.NumCols() {
		panic("Error: dataframes given to CombineDataFrames() must have same number of columns.")
	}

	newdf := append(df1, df2[1:]...)
	return newdf
}

/*
---------------------------------
Exploratory Data Analysis - Numerical version (COMPLETE)
---------------------------------
*/
//FiveNumSummary is a function that takes in an array of floats and returns a five-long array of floats
//corresponding to the minimum, 1st quartile, median, 3rd quartile, and maximum.
func FiveNumSummary(dataset []float64) []float64 {
	results := make([]float64, 5)
	data := Sort(dataset)
	n := len(data)

	//statistical function calls
	results[0] = data[0]
	results[1] = Median(data[0:n/2]) //1st quartile
	results[2] = Median(data)
	results[3] = Median(data[n/2:n]) //3rd quartile
	results[4] = data[n-1]
	
	return results
}

//Mean is a function that takes in an array of floats and returns the mean. Simple as that!
func Mean(data []float64) float64 {
	n := len(data)
	if n <= 0 {
		panic("Invalid dataset given to Mean().")
	}

	avg := 0.0
	for i := range data {
		avg += data[i] / float64(n)
	}

	return avg
}

//RecursiveMean is a function that takes in an array of floats and returns the mean, but recursively!
//This is a good test or practice question, idk.
func RecursiveMean(data []float64) float64 {
	n := float64(len(data))

	//error handling: empty dataset, or negative number of elements (much bigger problem)
	if n <= 0.0 {
		panic("Invalid dataset given to RecursiveMean().")
	} else if n == 1.0 { //base case: one element
		return data[0]
	} else { //recursive case: multiple elements
		//constants for help later
		c1 := (n - 1.0)/n
		c2 := 1.0/n

		currVal := data[0]
		recMean := RecursiveMean(data[1:])
		return c1 * recMean + c2 * currVal
	}
}

//Sort is a function that takes in an array of floats and returns the sorted array via selection sort.
func Sort(data []float64) []float64 {
	if len(data) == 1 {
        return data
    }
    
    newList := make([]float64, 0)
    firstIndex := len(data) - 1
    smallestVal := data[0]
    for i := range data {
        if data[i] < smallestVal {
            smallestVal = data[i]
        }
    }
    
    //now that we have smallest value, range through again and sort out everything else
    for i := range data {
        if data[i] == smallestVal {
            if i < firstIndex {
                firstIndex = i
            }
        }
    }
    
    newList = append(newList, data[:firstIndex]...)
    newList = append(newList, data[firstIndex+1:]...)
    
    sortedList := make([]float64, 1)
    sortedList[0] = smallestVal
    sortedRest := Sort(newList)
    
    return append(sortedList, sortedRest...)
}

//Median is a function that takes in an array of floats, sorts it, and returns the median.
func Median(dataset []float64) float64 {
	data := Sort(dataset)
	n := len(data)

	medIndex := (n/2)
	if n % 2 == 1 { //n is odd
		
		return data[medIndex]
	}
	//else it is even
	medIndex1 := n/2 - 1

	med := (data[medIndex] + data[medIndex1]) / 2
	return med
}

//Variance is a function that takes in an array of floats and returns the variance.
func Variance(dataset []float64) float64 {
	n := float64(len(dataset))
	xMean := Mean(dataset)

	//loop through entries, calculations
	VarX := 0.0
	for i := range dataset {
		currVal := (dataset[i] - xMean) * (dataset[i] - xMean)
		VarX += currVal
	}

	VarX = VarX / (n - 1.0)
	return VarX
}

//StdDev is a function that takes the square root of the variance. Nothing to it lmao
func StdDev(dataset []float64) float64 {
	return math.Sqrt(Variance(dataset))
}

//Correlation takes two datasets and returns the Pearson's correlation coefficient between the two.
//Critically, it assumes that the two datasets are the same length, and will simply take numbers off the smaller one.
func Correlation(x, y []float64) float64 {
	//three "main" terms
	num := 0.0
	denom1 := 0.0
	denom2 := 0.0

	//error handling
	n1 := len(x)
	n2 := len(y)
	if n1 != n2 {
		fmt.Println("    Warning: Variables of unequal size given to Correlation(). Excess terms removed from the longer vector.")
	}

	n := 0
	if n1 > n2 {
		n = n2
	} else {
		n = n1
	}

	xMean := Mean(x)
	yMean := Mean(y)

	for i := 0; i < n; i++ {
		currX := (x[i] - xMean)
		currY := (y[i] - yMean)

		//adding to terms
		num += currX * currY
		denom1 += currX * currX
		denom2 += currY * currY
	}

	corr := num / (math.Sqrt(denom1) * math.Sqrt(denom2))
	return corr
}

//CorrMatrix takes in a dataframe and returns the correlation of all numeric variables with each other, as well
//as the names of the variables for each row/col.
func CorrMatrix(df dataframe) [][]float64 {
	//actual matrix
	var corrMat [][]float64
	for i := 0; i < df.NumCols(); i++ {
		var newRow []float64

		currentDF1 := df.RemoveNAs(i)
		col1 := currentDF1.GetColumn(i)[1:]

		for j := 0; j < df.NumCols(); j++ {
			currentDF2 := df.RemoveNAs(j)
			col2 := currentDF2.GetColumn(j)[1:]

			if IsNumeric(col1) && IsNumeric(col2) {
				corrVal := Correlation(MakeNumeric(col1), MakeNumeric(col2))
				newRow = append(newRow, corrVal)
			}
		}

		//one final check
		if IsNumeric(col1) {
			corrMat = append(corrMat, newRow)
		}
	}

	return corrMat
}

//GetNamesForCorrMat takes in a dataframe and returns all the labels for the columns that are numeric.
func GetNamesForCorrMat(df dataframe) []string {
	var names []string
	
	//range through and get the names
	for i := range df[0] {
		currentDF := df.RemoveNAs(i)
		currentCol := currentDF.GetColumn(i)[1:]

		if IsNumeric(currentCol) {
			names = append(names, df[0][i])
		}
	}

	return names
}

//GetSummaryStats takes in a column of floats and returns an array of floats of the following:
//the mean, variance, standard deviation, and 5-number summary for the input data.
func GetSummaryStats(data []float64) []float64 {
	var sumStats []float64
	sumStats = append(sumStats, Mean(data))
	sumStats = append(sumStats, Variance(data))
	sumStats = append(sumStats, StdDev(data))

	sumStats = append(sumStats, FiveNumSummary(data)...)

	return sumStats
}

//GetSummaryStatsTableRow takes in a column of floats with a name and adds the name to the 
//beginning of the list, converting it to a string array.
func GetSummaryStatsTableRow(data []float64, name string) []string {
	stats := GetSummaryStats(data)

	var tabRow []string
	tabRow = append(tabRow, name)
	for i := range stats {
		val := strconv.FormatFloat(stats[i], 'f', -1, 64)
		tabRow = append(tabRow, val)
	}

	return tabRow
}

//GetSummaryStatsTable takes in a dataframe and returns a matrix of the summary statistics for
//all numerical variables in the dataframe.
func GetSummaryStatsTable(df dataframe) dataframe {
	var sumStatTab dataframe
	heading := make([][]string, 1)
	heading[0] = make([]string, 9)
	heading[0][0] = "variable"
	heading[0][1] = "mean"
	heading[0][2] = "variance"
	heading[0][3] = "stdev"
	heading[0][4] = "min"
	heading[0][5] = "1stQuartile"
	heading[0][6] = "median"
	heading[0][7] = "3rdQuartile"
	heading[0][8] = "max"

	//range through cols exclusively
	for i := range df[0] {
		currentDF := df.RemoveNAs(i)
		currentCol := currentDF.GetColumn(i)
		name :=  currentCol[0]

		if IsNumeric(currentCol[1:]) {
			newRow := GetSummaryStatsTableRow(MakeNumeric(currentCol[1:]), name)
			sumStatTab = append(sumStatTab, newRow)
		}
	}

	finalStats := append(heading, sumStatTab...)
	return finalStats
}

/*----------------
MISCELLANEOUS
---------------*/

//PrintInstructions forms the set of print statements that are printed to the console in case of a commandline argument set
//that does not match with what is allowed.
func PrintInstructions() {
	fmt.Println("This is the Golang portion of the project Invasive Species Modeling by Jonathan Zhu for Programming for Scientists 02-601.")
	fmt.Println("")
	fmt.Println("To begin, you need a csv file of occurrence data for some species from the Global Biodiversity Information Facility.")
	fmt.Println("    Some are provided in the same directory as this file.")
	fmt.Println("Once you have the raw data, make sure it does not have any erroneous commas and is formatted correctly. This can be done")
	fmt.Println("    by using a spreadsheet tool like Numbers or Excel. More info is given in README.md.")
	fmt.Println("")
	fmt.Println("To run this program, after typing ./project into the command line, add the following arguments:")
	fmt.Println("1. The mode to run, either `preprocessing` or `eda`. Preprocessing mode will format the data correctly and put it into a file called out.csv.")
	fmt.Println("   EDA mode will compute some summary statistics and correlation calculations on the formatted data.")
	fmt.Println("2. The file name to be read.")
	fmt.Println("3. If running preprocessing, the number of pseudo-absences to generate. Leave blank for the default, which is twice the number of presences.")
	fmt.Println("")
	fmt.Println("There are helpful comments included throughout the files as you journey as well. Good luck!")
}