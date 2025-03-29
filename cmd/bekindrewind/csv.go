package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
)

func csvToMovies(file multipart.File) ([]*Movie, error) {
	// read the CSV file but skip the first line
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	records = records[1:]

	// parse the CSV records into Movie structs
	movies := make([]*Movie, len(records))
	for i, record := range records {
		year, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("error parsing csv at line %d, invalid year: %w", i+1, err)
		}
		seen, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, fmt.Errorf("error parsing csv at line %d, invalid seen_year: %w", i+1, err)
		}
		rate, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, fmt.Errorf("error parsing csv at line %d, invalid rate: %w", i+1, err)
		}
		movies[i] = &Movie{
			Title:     record[1],
			Sagas:     strings.Split(record[2], "|"),
			Year:      uint16(year),
			SeenYear:  uint16(seen),
			Rate:      uint8(rate),
			Directors: strings.Split(record[6], "|"),
			Writers:   strings.Split(record[7], "|"),
			Composers: strings.Split(record[8], "|"),
			Dops:      strings.Split(record[9], "|"),
			Editors:   strings.Split(record[10], "|"),
			Producers: strings.Split(record[11], "|"),
			Studios:   strings.Split(record[12], "|"),
			Countries: strings.Split(record[13], "|"),
			Genres:    strings.Split(record[14], "|"),
		}
	}

	Logger.Infof("Loaded %d movies from file %+v", len(records), file)
	return movies, nil
}

func moviesToCSV(movies []*Movie) (string, error) {
	// open the CSV file
	buffer := io.Writer(&strings.Builder{})

	// create a new CSV writer
	writer := csv.NewWriter(buffer)
	defer writer.Flush()

	// write the header
	if err := writer.Write([]string{
		"id",
		"title",
		"saga",
		"release_year",
		"seen_year",
		"vote",
		"directors",
		"writers",
		"composers",
		"dops",
		"editors",
		"producers",
		"studios",
		"countries",
		"genres",
	}); err != nil {
		return "", fmt.Errorf("error writing the header: %w", err)
	}

	// write the records
	Logger.Info("writing records")
	for _, movie := range movies {
		if err := writer.Write([]string{
			strconv.Itoa(movie.Id),
			movie.Title,
			strings.Join(movie.Sagas, "|"),
			strconv.Itoa(int(movie.Year)),
			strconv.Itoa(int(movie.SeenYear)),
			strconv.Itoa(int(movie.Rate)),
			strings.Join(movie.Directors, "|"),
			strings.Join(movie.Writers, "|"),
			strings.Join(movie.Composers, "|"),
			strings.Join(movie.Dops, "|"),
			strings.Join(movie.Editors, "|"),
			strings.Join(movie.Producers, "|"),
			strings.Join(movie.Studios, "|"),
			strings.Join(movie.Countries, "|"),
			strings.Join(movie.Genres, "|"),
		}); err != nil {
			return "", fmt.Errorf("error writing the record %v: %w", movie, err)
		}
	}

	Logger.Infof("wrote %d movies", len(movies))
	return buffer.(*strings.Builder).String(), nil
}
