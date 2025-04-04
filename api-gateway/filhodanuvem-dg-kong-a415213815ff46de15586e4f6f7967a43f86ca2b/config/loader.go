package config

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/devgymbr/kong"
)

func Loader(config *kong.Config, fileName string, tickerDuration time.Duration) error {
	fp, err := openFile(fileName)
	if err != nil {
		return err
	}

	if err = refreshConfig(config, fp); err != nil {
		return err
	}
	fp.Close()

	go func() {
		defer fp.Close()
		t := time.NewTicker(tickerDuration)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				slog.Debug("checking config file for changes")
				fp, err = openFile(fileName)
				if err != nil {
					slog.Error("could not open config file: %s", err)
					continue
				}

				if err = refreshConfig(config, fp); err != nil {
					slog.Error("could not refresh config: %s", err)
					fp.Close()
					continue
				}
				fp.Close()
			}
		}
	}()

	return nil
}

func openFile(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0755)
}

func refreshConfig(config *kong.Config, fp *os.File) error {
	data, err := io.ReadAll(fp)
	if err != nil {
		return err
	}

	info, err := fp.Stat()
	if err != nil {
		return err
	}

	modTime := info.ModTime()
	if config.ModifiedSince(modTime) {
		return nil
	}

	slog.Debug("updating config based on changes")

	err = config.Refresh(data, modTime)

	slog.Info("config updated")

	return err
}
