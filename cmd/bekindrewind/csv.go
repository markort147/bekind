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

		rolePeople := make(map[*string][]string)
		for index, role := range map[int]string{
			6:  "director",
			7:  "writer",
			8:  "composer",
			9:  "cinematographer",
			10: "editor",
			11: "producer",
		} {
			if record[index] != "" {
				for _, name := range strings.Split(record[index], "|") {
					if rolePeople[&name] == nil {
						rolePeople[&name] = make([]string, 0)
					}
					rolePeople[&name] = append(rolePeople[&name], role)
				}
			}
		}

		movies[i] = &Movie{
			Title:     record[1],
			Sagas:     strings.Split(record[2], "|"),
			Year:      year,
			SeenYear:  seen,
			Rate:      rate,
			People:    rolePeople,
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

		directors := make([]string, 0)
		writers := make([]string, 0)
		composers := make([]string, 0)
		dops := make([]string, 0)
		editors := make([]string, 0)
		producers := make([]string, 0)
		for name, roles := range movie.People {
			for _, role := range roles {
				switch role {
				case "director":
					directors = append(directors, *name)
				case "writer":
					writers = append(writers, *name)
				case "composer":
					composers = append(composers, *name)
				case "cinematographer":
					dops = append(dops, *name)
				case "editor":
					editors = append(editors, *name)
				case "producer":
					producers = append(producers, *name)
				default:
					return "", fmt.Errorf("unknown person role %q", role)
				}
			}
		}

		if err := writer.Write([]string{
			strconv.Itoa(movie.Id),
			movie.Title,
			strings.Join(movie.Sagas, "|"),
			strconv.Itoa(movie.Year),
			strconv.Itoa(movie.SeenYear),
			strconv.Itoa(movie.Rate),
			strings.Join(directors, "|"),
			strings.Join(writers, "|"),
			strings.Join(composers, "|"),
			strings.Join(dops, "|"),
			strings.Join(editors, "|"),
			strings.Join(producers, "|"),
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
