package fs

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FSDataSource struct {
	root string
}

type fsDataSourceJSON struct {
	Root string "json:root"
}

func NewFSDataSource(root string) (FSDataSource, error) {
	return FSDataSource{
		root: root,
	}, nil
}

func (this *FSDataSource) GetType() string {
	return "fs"
}

func (this *FSDataSource) GetReader() (io.Reader, error) {
	reader, writer := io.Pipe()
	go runTar(this.root, writer) // Ignore errors until we figure out how to propagate
	return reader, nil

}

func runTar(src string, writer *io.PipeWriter) {
	defer writer.Close()
	err := tarDir(src, writer)
	if err != nil {
		fmt.Printf("Err returned from tarDir %s\n", err.Error())
	} else {
		fmt.Printf("tarDir exited successfully\n")
	}
}

// Tar takes a source and variable writers and walks 'source' writing each file
// found to the tar writer
func tarDir(src string, writer io.Writer) error {

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("unable to tar files - %v", err.Error())
	}

	tw := tar.NewWriter(writer)
	defer tw.Close()
	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		fmt.Printf("walk file = %s\n", file)
		// return on any error
		if err != nil {
			return err
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			fmt.Printf("Skipping file = %s, not a regular file\n", file)

			return nil
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.CopyBuffer(tw, f, make([]byte, 1024*1024)); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()
		fmt.Printf("Finished writing file %s\n", file)
		return nil
	})
}

func (this FSDataSource) MarshalJSON() ([]byte, error) {

	jsonStruct := fsDataSourceJSON{
		Root: this.root,
	}

	return json.Marshal(jsonStruct)
}

func (this *FSDataSource) UnmarshalJSON(data []byte) error {
	jsonStruct := fsDataSourceJSON{}
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	this.root = jsonStruct.Root
	return nil
}
