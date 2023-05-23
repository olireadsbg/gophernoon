# Gophernoon - Building A CSV Parser

## Project Requirements

* Download the sample CSV file containing organisations
    https://media.githubusercontent.com/media/datablist/sample-csv-files/main/files/organizations/organizations-1000.csv
* Open and read the file, using the builtin csv library, to parse the records to 
    structs
    * https://pkg.go.dev/os#Open
        * Don't forget to close the file when you're done with it.
        * https://go.dev/tour/flowcontrol/12
    * https://pkg.go.dev/encoding/csv
* Sort the data by a the `Name` field in alphabetical ascending order
    * https://pkg.go.dev/sort
    * https://gobyexample.com/sorting-by-functions
* Output the result to `std.out`
    * https://pkg.go.dev/fmt
    * https://gobyexample.com/hello-world

## Optional Tasks

Optional tasks can be picked up in any order. They have a difficulty (1-5) 
associated with them to help you decide what level of challenge you want.

### Make your application configurable (Difficulty: 1)

* Use a `input.file` flag to pass in a csv file when running the command
    ``` sh
    go build -o csvparser && ./csvparser -input.file ./data.csv
    go build -o csvparser && ./csvparser | cat ./data.csv
    ```
    * https://pkg.go.dev/flag
    * https://gobyexample.com/command-line-flags

### Read your CSV file line by line (Difficulty: 3)

* Your CSV file may be very large, and your system may not have much memory. 
    Limit the amount of memory to your Go application and load in the large CSV 
    file, parsing it line by line rather than loading the whole file into memory.
    * https://github.com/datablist/sample-csv-files/raw/main/files/organizations/organizations-2000000.zip
    * https://weaviate.io/blog/gomemlimit-a-game-changer-for-high-memory-applications

### Use Concurrency to implement a "Merge Sort" (Difficulty: 5)

* We can optimise the sorting of a large data set using a recursive split and 
    merge method. Use waitgroups and go routines to achieve this. Be sure to use
    contexts too incase of any parsing erros.
    * http://pages.di.unipi.it/marino/pythonads/SortSearch/TheMergeSort.html
    * https://go.dev/tour/concurrency/1
    * https://gobyexample.com/waitgroups
    * https://pkg.go.dev/context
    * https://gobyexample.com/context