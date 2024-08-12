package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	//there are two main parts to this: the preprocessing, and the EDA.
	//If we get the argument "preprocessing", we execute a preprocessing pipeline, described right below.
		//Here's how this general pipeline works: 
		//we take a CSV of occurrences from GBIF (Global biodiversity Information Facility)
		//The filename for this is given in the second argument.
		//these are free to access and are (generally) in the same format for all of them

		//So we will remove certain columns outright; all columns that are not basic identifiers (scientific name)
		//critically, we DO NOT remove the data for locations (lat/long/elevation/depth).
		//We then remove all the NA values for lat and long; we cannot really work without these.

		//After this, we get the random points. We add AT MOST three times as many points as we have remaining 
		//in the dataset, or we take in from the user. We then filter them.

		//We then write the CSV file to be worked on in R. 
	
	//now if we get the argument "eda", we execute functions that give us numerical EDA for a numerical dataset. 
		//Here's how this general pipeline works:
		//we take a numerical CSV outputted by the first R file, and then simply put it into the functions for summary statistics
		//and correlation matrix. 

		//We write these summary statistics and matrix to .txt files. We can visualize them in R, if we want. 

	//Now begins the actual code for this pipeline.

	//first, we'll write an error handling.
	if len(os.Args) == 1 {
		PrintInstructions()
		os.Exit(2)
	}

	mode := os.Args[1]
	if mode == "preprocessing" {
		fmt.Println("Preprocessing pipeline selected.")

		//read the file
		df := ReadCSV(os.Args[2])

		fmt.Println("File read successfully. Now processing")

		//remove the columns we want to remove
		//to do this, we'll use a few loops
		for i := 0; i < 11; i++ {
			df = df.DropColumn(1)

			/*here's what this loop removes, in order:
			1. datasetKey
			2. occurrenceID
			3. kingdom
			4. phylum
			5. class
			6. order
			7. family
			8. genus
			9. species
			10. infraspecificEpithet
			11. taxonRank

			*/
		}

		for i := 0; i < 8; i++ {
			df = df.DropColumn(2)

			/* here's what this loop removes, in order:
			1. verbatimScientificName
			2. verbatimScientificNameAuthorship
			3. countryCode
			4. locality
			5. stateProvince
			6. occurrenceStatus
				note: we remove this because it is listed as "PRESENT", a string, in the CSV, and we want numeric.
			7. individualCount
			8. publishingOrgKey

			*/
		}

		for i := 0; i < 2; i++ {
			df = df.DropColumn(4)

			//removes coordinateUncertaintyInMeters and coordinatePrecision
		}

		df = df.DropColumn(5) //removes elevationAccuracy

		for i := 0; i < 22; i++ {
			df = df.DropColumn(6) //should remove the rest of the columns
		}

		//now we remove the NA values in rows where latitude/longitude is NA.
		//I also did for elevation.
		//elevation and depth can be imputed later, but these ones are crucial.
		df = df.RemoveNAs(2)
		df = df.RemoveNAs(3)
		df = df.AddColOfOnes("present")
		
		fmt.Println("Processing complete. Adding Random Points")
		//if you want, you can specify a random number of points. 
		//but you don't have to, and the below code handles any errors that might occur if you don't.
		n := len(os.Args)
		if n == 3 { //pass in number that's two times number of rows of current dataframe
			n = 2 * df.NumRows()
			pseudoAbs := GetRandPoints(n)

			//dataframe with the pseudoabsences
			df_pas := SortRandPoints(df, pseudoAbs)

			fmt.Println("Random points added. Writing CSV")
			WriteCSV(df_pas, "out.csv")
		} else { //we got a number as the third arg, so we pass that
			numPoints, err1 := strconv.Atoi(os.Args[3])
			if err1 != nil {
				panic(err1)
			}

			pseudoAbs := GetRandPoints(numPoints)

			//dataframe with the pseudoabsences
			df_pas := SortRandPoints(df, pseudoAbs)

			fmt.Println("Random points added. Writing CSV")
			WriteCSV(df_pas, "out.csv")
		}
	
	fmt.Println("Processed data written to out.csv. Use it in the next file, ml_process.rmd.")
	//END OF PREPROCESSING MODULE

	} else if mode == "eda" {
		fmt.Println("EDA process selected")

		//writing the summary statistics
		fmt.Println("Writing Summary Statistics CSV")
		df := ReadCSV(os.Args[2])
		sumStats := GetSummaryStatsTable(df)
		WriteCSV(sumStats, "summary.csv")

		//writing corr matrix\
		fmt.Println("Writing Correlation Matrix")
		corrMat := CorrMatrix(df)
		names := GetNamesForCorrMat(df)
		WriteCorrMatrixToFile(corrMat, names, "corrmat.csv")
		
		fmt.Println("Summary statistics written to summary.csv, and correlation matrix written to corrmat.csv.")
		fmt.Println("Use data_full.csv for the next file, ml_models.rmd.")

	} else { //error handling in case of other arguments given: give instructions
		PrintInstructions()
	}
}