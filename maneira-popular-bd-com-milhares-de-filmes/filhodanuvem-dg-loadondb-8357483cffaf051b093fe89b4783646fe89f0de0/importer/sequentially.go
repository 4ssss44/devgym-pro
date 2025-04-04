package importer

import (
	"context"
	"encoding/csv"
	"io"
	"log"

	loadondb "github.com/filhodanuvem/dg-loadondb"
	"github.com/filhodanuvem/dg-loadondb/database"
	"github.com/filhodanuvem/dg-loadondb/errors"
	"github.com/filhodanuvem/dg-loadondb/parser"
)

type Sequentially struct {
	Writer    database.Writer
	ChunkSize int
}

func (s Sequentially) ImportMovies(ctx context.Context, reader *csv.Reader) error {
	log.Printf("Inserting movies sequentially with chunckSize=%d", s.ChunkSize)
	i := -1
	movies := []loadondb.Movie{}
	errors := errors.List{}
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if i == -1 {
			i++
			continue
		}
		movie, err := parser.ParseLine(line)
		if err != nil {
			errors.AddError(err)
			continue
		}
		movies = append(movies, movie)
		i++
		if i%(s.ChunkSize) == 0 {
			if err = s.Writer.InsertMovies(ctx, movies); err != nil {
				errors.AddError(err)
			}
			movies = []loadondb.Movie{}
		}
	}
	if len(movies) > 0 {
		if err := s.Writer.InsertMovies(ctx, movies); err != nil {
			errors.AddError(err)
		}
	}

	if errors.Len() == 0 {
		return nil
	}

	return &errors
}
