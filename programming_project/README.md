# Invasive Species Modeling
## Jonathan Zhu, Programming for Scientists 02-601

This is the documentation that is (hopefully) sufficient for you to use all the files necessary in this project.
Included are the following files:
- All the necessary Golang files, including functions.go, functions_test.go, and main.go. 
- All the necessary R files, which should be ml_process.rmd, ml_models.rmd, and ml_biomod.rmd.
- A few datasets in the /data/ folder for you to get started.
- The necessary climactic data, stored in the folder /wc10/.

To begin with, make sure you have the following Golang packages installed:
- `fmt`
- `os`
- `bufio`
- `strings`
- `strconv`
- `math`
- `math/rand`
- `io/fs`
- `testing` (for testing purposes, optional if you don't want to test my code :D)

Additionally, make sure you have the following R packages installed:
- `tidyverse`
- `tidymodels`
- `readr`
- `raster`
- `dplyr` (included in tidymodels)
- `data.table`
- `rgbif`
- `maps`
- `biomod2`
- `tidyterra`
- `ggtext`

The code included in this bundle will help you generate a projection of invasive species spread from datafile to plot. First, you will need some data. All the data I used in this project is from the Great Biodiversity Information Facility (GBIF) (gbif.org), from which CSV files of occurrences can be downloaded for free and for which this project specifically caters to. There is some sample data located in the /data/ folder formatted correctly. IMPORTANT: If you get data straight from GBIF, it will not be directly usable in Golang. You will need to open it in a spreadsheet editor, like Excel or Numbers (which is the default on Macs), replace every comma with a semicolon (or other non-comma character) and export this new file as a CSV. Unfortunately, my code is not sophisticated enough to detect commas that are within entries and not separating them.

Once you're ready to start, first, you will need to run the Go pipeline, so navigate to the directory with main.go and functions.go and in the commandline type "./project preprocessing [your file name]", obviously using your filename as the last argument. This will generate a CSV file with a predetermined number of pseudoabsences in the world, though you can add this as a third argument.

Next, open ml_process.rmd. Make sure your file from the previous section is labeled "out.csv", which it is by default. Run all the code, and you will get a preprocessed data file called "data_full.csv". Following this, it is optional but you can use this in golang and type in the commandline "./project eda data_full.csv". This will give you a CSV of numerical summary called summary.csv and a correlation matrix of correlation values called corrmat.csv.

The main machine learning component of this project is then found in ml_models.rmd. Run all the code with the libraries, the recipe, and the datasplit, which are the first three code chunks. Following this, the document is then labeled so you can run GLM, Random Forest, or XGBoost on your data with regression, or Random Forest or XGBoost for classification. Pick your choice of which to run; make sure that if you run XGBoost, that you have a BEEFY computer! Once that is done, you can use the ending code for plotting purposes, which gives you the predicted spread plots and the actual occurrence plot. Run these selectively based on the model you ran. Alternatively, you can run everything, which takes a while but is what I did for my initial run with the Parthenium weed occurrences.

Finally, if you wish to compare the models produced here to the main standard in spread plotting, open ml_biomod.rmd and make sure the datafile name matches your data file. Run all the code and you will get the plots for predicted spread via GLM, Random forest, and XGBoost, as well as some metrics for the various runs it does. Again, this requires a pretty beefy computer.

If you want to test my functions in Golang, type "go test -v" into the command line.

That's just about all you need for now. For any questions or concerns, please contact me at jazhu@andrew.cmu.edu.
Thanks!
  Jonathan

P.S.: I know I didn't get to 2000 lines of code on the Golang portion of this project, but I reached 2000 with the help of my R code. Please forgive me for that, but I felt as though many of my functions were too big to have test cases for. That and I didn't need a huge number of functions in Golang for this project, because a lot of the heavy lifting (machine learning) was done in R.