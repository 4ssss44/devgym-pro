package parser

import (
	"regexp"
	"strconv"
	"strings"

	loadondb "github.com/filhodanuvem/dg-loadondb"
	"github.com/filhodanuvem/dg-loadondb/errors"
)

func ParseLine(line []string) (loadondb.Movie, error) {
	if len(line) < 3 {
		return loadondb.Movie{}, errors.NewNonFormattedLine(line)
	}

	id, err := strconv.Atoi(line[0])
	if err != nil {
		return loadondb.Movie{}, errors.NewNonValidID(line[0])
	}

	title := line[1]
	year := 0
	re, err := regexp.Compile("(.*)\\s*\\((.*)\\)")
	if err == nil {
		tileMatches := re.FindStringSubmatch(line[1])
		if len(tileMatches) == 3 {
			title = strings.Trim(tileMatches[1], " ")
			year, err = strconv.Atoi(tileMatches[2])
			if err != nil {
				// nothing, year as numeric is optional
			}
		}
	}

	genres := strings.Split(strings.Trim(line[2], "\""), "|")
	if len(genres) == 1 && genres[0] == "" {
		genres = []string{}
	}
	return loadondb.Movie{
		ID:     id,
		Title:  title,
		Year:   year,
		Genres: genres,
	}, nil
}
