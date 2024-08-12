package main

import (
	"bufio"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
	"math"
)

//testing functions!
//so I'm going to be using some of the premade CSV files for some of these. 

/* -------------
TYPES
----------------*/

// Here, I am aiming to test all functions that do not have to do with in/out files.

//A lot of my testing functions will be for the EDA portion of this project.
//Many of my functions work with large data files (CSVs) and I don't want to take up so much space that
//I won't be able to submit the file.

type MedianTest struct {
	dataset []float64
	result float64
}

type MeanTest struct {
	dataset []float64
	result float64
}

type VarianceTest struct {
	dataset []float64
	result float64
}

type StdDevTest struct {
	dataset []float64
	result float64
}

type CorrelationTest struct {
	x, y []float64
	result float64
}

type FiveNumSummaryTest struct {
	dataset []float64
	result []float64
}

type GetSummaryStatsTest struct {
	dataset []float64
	result []float64
}
/*---------------
Actual Test Functions: EDA
As well as specific readin functions!
----------------*/

func TestMedian(t *testing.T) {
	//read in all tests from the Tests/Richness directory and run them
	tests := ReadMedianTests("Tests/Median/")
	for _, test := range tests {
		//run the test
		result := Median(test.dataset)
		//check the result
		if roundFloat(result, 4) != test.result {
			t.Errorf("Median(%v) = %v, want %v", test.dataset, result, test.result)
		}
	}
}

func ReadMedianTests(directory string) []MedianTest {

	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]MedianTest, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].dataset = ReadFloatArrayFromFile(directory + "/input/" + inputFile.Name())
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadFloatFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

func TestMean(t *testing.T) {
	//read in all tests from the Tests/Richness directory and run them
	tests := ReadMeanTests("Tests/Mean/")
	for _, test := range tests {
		//run the test
		result := Mean(test.dataset)
		//check the result
		if roundFloat(result, 4) != test.result {
			t.Errorf("Mean(%v) = %v, want %v", test.dataset, result, test.result)
		}
	}
}

func TestRecursiveMean(t *testing.T) {
	//read in all tests from the Tests/Richness directory and run them
	tests := ReadMeanTests("Tests/Mean/")
	for _, test := range tests {
		//run the test
		result := RecursiveMean(test.dataset)
		//check the result
		if roundFloat(result, 4) != test.result {
			t.Errorf("Mean(%v) = %v, want %v", test.dataset, result, test.result)
		}
	}
}

func ReadMeanTests(directory string) []MeanTest {

	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]MeanTest, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].dataset = ReadFloatArrayFromFile(directory + "/input/" + inputFile.Name())
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadFloatFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

func TestVariance(t *testing.T) {
	//read in all tests from the Tests/Richness directory and run them
	tests := ReadVarianceTests("Tests/Variance/")
	for _, test := range tests {
		//run the test
		result := Variance(test.dataset)
		//check the result
		if roundFloat(result, 4) != test.result {
			t.Errorf("Variance(%v) = %v, want %v", test.dataset, result, test.result)
		}
	}
}

func ReadVarianceTests(directory string) []VarianceTest {

	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]VarianceTest, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].dataset = ReadFloatArrayFromFile(directory + "/input/" + inputFile.Name())
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadFloatFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

func TestStdDev(t *testing.T) {
	//read in all tests from the Tests/Richness directory and run them
	tests := ReadStdDevTests("Tests/StdDev/")
	for _, test := range tests {
		//run the test
		result := StdDev(test.dataset)
		//check the result
		if roundFloat(result, 4) != test.result {
			t.Errorf("StdDev(%v) = %v, want %v", test.dataset, result, test.result)
		}
	}
}

func ReadStdDevTests(directory string) []StdDevTest {

	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]StdDevTest, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].dataset = ReadFloatArrayFromFile(directory + "/input/" + inputFile.Name())
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadFloatFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

func TestFiveNumSummary(t *testing.T) {
	//read in all tests from the Tests/Richness directory and run them
	tests := ReadFiveNumSummaryTests("Tests/FiveNumSummary/")
	for _, test := range tests {
		//run the test
		result := FiveNumSummary(test.dataset)
		//check the result
		if !AreFloatArraysEqual(result, test.result) {
			t.Errorf("FiveNumSummary(%v) = %v, want %v", test.dataset, result, test.result)
		}
	}
}

func ReadFiveNumSummaryTests(directory string) []FiveNumSummaryTest {

	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]FiveNumSummaryTest, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].dataset = ReadFloatArrayFromFile(directory + "/input/" + inputFile.Name())
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadFloatArrayFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

func TestGetSummaryStats(t *testing.T) {
	//read in all tests from the Tests/Richness directory and run them
	tests := ReadGetSummaryStatsTests("Tests/GetSummaryStats/")
	for _, test := range tests {
		//run the test
		result := GetSummaryStats(test.dataset)
		//check the result
		if !AreFloatArraysEqual(result, test.result) {
			t.Errorf("GetSummaryStats(%v) = %v, want %v", test.dataset, result, test.result)
		}
	}
}

func ReadGetSummaryStatsTests(directory string) []GetSummaryStatsTest {

	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]GetSummaryStatsTest, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].dataset = ReadFloatArrayFromFile(directory + "/input/" + inputFile.Name())
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadFloatArrayFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}


/*--------------
General In/Out Functions!!!
-----------------*/

// ReadDirectory reads in a directory and returns a slice of fs.DirEntry objects containing file info for the directory
func ReadDirectory(dir string) []fs.DirEntry {
	//read in all files in the given directory
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return files
}

//ReadDataframe takes in a string name that will be the filename for the file we want to read. 
//It returns a dataframe where each sub-array in the array corresponds to a row of the file.
//Critically, all values are separated by commas, much like a CSV.
func ReadDataframe(filename string) dataframe {
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

//ReadDataFrame takes in a string name that will be the filename for the CSV we want to read. 
//It returns a dataframe where each sub-array in the array corresponds to a row of the CSV.
func ReadDataFrame(filename string) dataframe {
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

//ReadDataframeWithSpaces takes in a string name that will be the filename for the file we want to read. 
//It returns a dataframe where each sub-array in the array corresponds to a row of the file.
//Critically, all values are separated by spaces, allowing us to have more types of files for test cases.
func ReadDataframeWithSpaces(filename string) dataframe {
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
		entries := strings.Split(currentRow, " ")

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

//ReadDataFrameWithNAs takes in a string name that will be the filename for the CSV we want to read. 
//It returns a dataframe where each sub-array in the array corresponds to a row of the CSV.
//Here, all NAs read are converted to "" the empty string.
func ReadDataFrameWithNAs(filename string) dataframe {
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
			if entries[i] == "NA" {
				df[currentIndex][i] = ""
			} else {
				df[currentIndex][i] = entries[i]
			}
			
		}

		currentIndex++
	}

	err2 := scanner.Err()
	if err2 != nil {
		panic(err2)
	}

	return df
}

//ReadBool reads in either "true" or "false" from a file and returns that as a boolean.
func ReadBool(file string) bool {
	var newBool bool

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "true" {
			newBool = true
		} else if line == "false" {
			newBool = false
		} else {
			panic("invalid boolean. Must read `true` or `false`.")
		}
	}

	return newBool
}

// ReadFloatFromFile reads in a single float from a file
func ReadFloatFromFile(file string) float64 {
	//open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//create a new scanner
	scanner := bufio.NewScanner(f)

	//read in the line
	scanner.Scan()
	line := scanner.Text()

	//convert the line to an int using strconv
	value, err := strconv.ParseFloat(line, 64)
	if err != nil {
		panic(err)
	}

	return value
}

// ReadFloatArrayFromFile reads in a single float array from a file
func ReadFloatArrayFromFile(file string) []float64 {
	//open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//create a new scanner
	scanner := bufio.NewScanner(f)

	//read in the line
	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, " ")

	var floatArray []float64
	for i := range parts {
		value, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			panic(err)
		}

		floatArray = append(floatArray, value)
	}

	//convert the line to an int using strconv
	

	return floatArray
}

// ReadTwoFloatArrays reads in two float arrays from a file.
//Critically, each one is on a new line.
//hypothetically, this could work for more than two arrays, but in the context of testing, we don't need more.
func ReadTwoFloatArrays(file string) [][]float64 {
	//open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var floatArrays [][]float64

	//create a new scanner
	scanner := bufio.NewScanner(f)

	//read in the line
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		var floatArray []float64
		for i := range parts {
			value, err := strconv.ParseFloat(parts[i], 64)
			if err != nil {
				panic(err)
			}

			floatArray = append(floatArray, value)
		}

		floatArrays = append(floatArrays, floatArray)
	}
	

	//convert the line to an int using strconv
	

	return floatArrays
}

/* ---------------
Equality and Rounding Functions
-----------------*/

//AreDFsEqual takes two dataframes df1, df2 and returns true if all elements in df1 are equal 
//to the corresponding elements in df2.
func AreDFsEqual(df1, df2 dataframe) bool {
	//first test dimensions. If rows and cols not equal, return false
	if df1.NumRows() != df2.NumRows() {
		return false
	}
	if df1.NumCols() != df2.NumCols() {
		return false
	} 

	//then test every element, since we know equal size.
	for i := range df1 {
		for j := range df1[i] {
			if df1[i][j] != df2[i][j] {
				return false
			}
		}
	}

	return true
}

//AreFloatArraysEqual takes two float arrays, a1 and a2. It returns true if all elements in a1 are equal to the corresponding
//elements in a2.
func AreFloatArraysEqual(a1, a2 []float64) bool {
	//range through and compare. simple as that!
	for i := range a1 {
		if roundFloat(a1[i], 4) != roundFloat(a2[i], 4) {
			return false
		}
	}

	//if we've gotten here, all values should be equal
	return true
}

// roundFloat rounds a float to a given number of decimals precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// ReadIntFromFile reads in int from a file
func ReadIntFromFile(file string) int {
	//open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//create a new scanner
	scanner := bufio.NewScanner(f)

	//read in the line
	scanner.Scan()
	line := scanner.Text()
	//convert the line to an int using strconv
	value, err := strconv.Atoi(line)
	if err != nil {
		panic(err)
	}

	return value
}