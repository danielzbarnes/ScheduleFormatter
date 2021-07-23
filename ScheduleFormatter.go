package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)


/*

Example Input Data:

695
8/30/21
Mon
The new series title, Pt. 1
The original series title (A Day)
reference
5/5/19
22:45

Example Output:

8/30/21
The new series title, Pt. 1
The new series title
Lance Quinn
reference
5/5/19


*/
func main() {

	// it's easier to see the text if stored as variable rather than doing the work all in the for loop
	var textLines []string

	file, err := os.Open("Schedule.txt") // schedule exported from excel

	// a bit of error checking
	if err != nil {
		fmt.Println("opening file error", err)
	}

	// scanner is used to read lines from the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan(){ // scan holds a boolean if there's text
		textLines = append(textLines, scanner.Text()) // text any text found
	}

	closeErr := file.Close()
	if closeErr != nil {
		return
	}

	var outputLines []string // array to hold only the data I want


	for _, val := range textLines {

		if _, err := strconv.Atoi(val); err != nil{ // exclude the program number

			if isNotDay(val) { // exlude the day

				if isNotTime(val){ // exlude the duration

					if isNotOriginalSeries(val){ // exclude the original series line

						outputLines = append(outputLines, val) // If not any of the above then add to the output array

						if strings.Contains(val, " Pt. "){
							outputLines = append(outputLines, titleToSeries(val))
							outputLines = append(outputLines, "Lance Quinn") // arbitrarily adding the author here
						}
					}
				}
			}
		}
	}

	outputErr := outputToFile(outputLines)

	if err != nil{
		log.Fatal(outputErr)
	}
}

// write to file
func outputToFile(output []string) error{

	file, err := os.Create("output.txt")

	if err != nil {
		fmt.Println("Error outputting file.", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	for _, val := range output {

		_, err := io.WriteString(file, val+"\n")
		if err != nil {
			return err
		}
	}
	return file.Sync()
}



// returns the title as series by spliting off the " Pt. #"
func titleToSeries(title string) string{

	// this could use some refinement to remove the comma from ", Pt. #", however title with "?", "!" need to be retained
	series := strings.Split(title, " Pt. ")

	return series[0]
}

// returns true if the line is not a 3 letter day
func isNotDay(line string) bool  {

	// length for day is always 3 letters
	if len(line) == 3 {

		// parsing through days of the week is just an added check in case the schedule format changes
		daysOfTheWeek := []string{"Mon", "Tue", "Wed", "Thu", "Fri"}

		for _, day := range daysOfTheWeek{
			if strings.Contains(line, day){
				return false
			}
		}
		return false // added check in case the for loop ends without returning false
	}

	return true
}

// returns true if the line is not the original series
func isNotOriginalSeries(line string) bool{
	
	if strings.Contains(line, "Day)") {
		return false
	}
	return true
}

// returns true if the line is not a timestamp
func isNotTime(line string) bool{

	// time will always have a 5 digits format "xx:xx"
	if len(line) == 5 {

		// remove the ":" to get just the numbers
		num := strings.Replace(line, ":", "", -1)

		// additional check to make sure there's only numbers
		if _, err := strconv.Atoi(num); err == nil{
			return false
		}

		return false // failsafe in case weird input gets through
	}
	return true
}
