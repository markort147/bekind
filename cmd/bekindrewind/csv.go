package main

import (
	"encoding/csv"
	mv "github.com/markort147/bekind/cmd/bekindrewind/pkg/movies"
	"github.com/markort147/gopkg/log"
	"io/fs"
	"os"
	"strconv"
)

func fromCSV(fileSystem fs.FS, filePath string) [][]string {
	file, err := fileSystem.Open(filePath)
	if err != nil {
		log.Logger.Fatal(err)
	}
	defer func(file fs.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Logger.Fatal(err)
	}

	//log.Logger.Debugf("Read %d records", len(records))
	// return all but the first row
	return records[1:]
}

func CSVToMovies(fileSystem fs.FS, filePath string) []int {
	records := fromCSV(fileSystem, filePath)
	ids := make([]int, len(records))
	for i, record := range records {
		// title is 1, year is 3, rating is 5
		title := record[1]
		year, _ := strconv.Atoi(record[3])
		rate, _ := strconv.Atoi(record[5])
		movie := mv.Movie{Title: title, Year: uint16(year), Rate: uint8(rate)}
		saved := mv.Save(movie)
		ids[i] = saved.Id
	}
	return ids
}

func MoviesToCSV(filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Logger.Fatal(err)
	}
	defer func(file fs.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	err = writer.Write([]string{"id", "title", "year", "rate"})
	if err != nil {
		log.Logger.Fatal(err)
	}

	for _, movie := range mv.Data.Movies {
		err = writer.Write([]string{strconv.Itoa(movie.Id), movie.Title, strconv.Itoa(int(movie.Year)), strconv.Itoa(int(movie.Rate))})
		if err != nil {
			log.Logger.Fatal(err)
		}
	}

}
