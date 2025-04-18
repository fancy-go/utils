package file_util

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

func TarGzDir(dir string, archive string) (err error) {
	out, err := os.Create(archive)
	if err != nil {
		return fmt.Errorf("os.Create(%s) failed: %v", archive, err)
	}
	defer func() {
		closeErr := out.Close()
		if err == nil {
			err = errors.Wrapf(closeErr, "close archive %s error", archive)
		} else if closeErr != nil {
			log.Printf("close archive %s error: %v", archive, closeErr)
		}
		if err != nil {
			removeErr := os.Remove(archive)
			if removeErr != nil {
				fmt.Printf("remove %s failed: %v", archive, removeErr)
			}
		}
	}()

	bufOut := bufio.NewWriter(out)
	// Create a new tar archive
	gz := gzip.NewWriter(bufOut)
	tw := tar.NewWriter(gz)
	base := filepath.Base(dir)
	err = filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "walk file error")
		}
		if fi.IsDir() {
			return nil
		}
		header, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return errors.Wrap(err, "file info header error")
		}

		header.Name, err = filepath.Rel(dir, file)
		if err != nil {
			return errors.Wrapf(err, "get relative path error")
		}
		header.Name = filepath.Join(base, header.Name)
		if err = tw.WriteHeader(header); err != nil {
			return errors.Wrap(err, "write header error")
		}
		fileReader, err := os.Open(file)
		if err != nil {
			return errors.Wrapf(err, "open file %s error", file)
		}
		defer fileReader.Close()
		if _, err = io.Copy(tw, fileReader); err != nil {
			return errors.Wrapf(err, "copy file %s error", file)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "tar file error")
	}

	if err = tw.Close(); err != nil {
		return errors.Wrap(err, "close tar writer error")
	}

	if err = gz.Close(); err != nil {
		return errors.Wrap(err, "close gzip writer error")
	}

	if err = bufOut.Flush(); err != nil {
		return errors.Wrap(err, "flush tar writer error")
	}

	return nil
}
