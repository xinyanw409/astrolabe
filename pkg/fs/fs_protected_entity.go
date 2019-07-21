package fs

import (
	"archive/tar"
	"fmt"
	"github.com/vmware/arachne/pkg/arachne"
	vim "github.com/vmware/govmomi/vim25/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	//	"github.com/vmware/govmomi/vslm"
	"context"
)

type FSProtectedEntity struct {
	fspetm   *FSProtectedEntityTypeManager
	id       arachne.ProtectedEntityID
	name     string
	root     string
	data     []arachne.DataTransport
	metadata []arachne.DataTransport
	combined []arachne.DataTransport
}

func newProtectedEntityID(id vim.ID) arachne.ProtectedEntityID {
	return arachne.NewProtectedEntityID("fs", id.Id)
}

func newFSProtectedEntity(fspetm *FSProtectedEntityTypeManager, id arachne.ProtectedEntityID,
	name string, root string) (FSProtectedEntity, error) {
	data, metadata, combined, err := fspetm.getDataTransports(id)
	if err != nil {
		return FSProtectedEntity{}, err
	}
	newFSPE := FSProtectedEntity{
		fspetm:   fspetm,
		id:       id,
		name:     name,
		root:     root,
		data:     data,
		metadata: metadata,
		combined: combined,
	}
	return newFSPE, nil
}
func (this FSProtectedEntity) GetInfo(ctx context.Context) (arachne.ProtectedEntityInfo, error) {
	retVal := arachne.NewProtectedEntityInfo(
		this.id,
		this.name,
		this.data,
		this.metadata,
		this.combined,
		[]arachne.ProtectedEntityID{})
	return retVal, nil
}

func (this FSProtectedEntity) GetCombinedInfo(ctx context.Context) ([]arachne.ProtectedEntityInfo, error) {
	fsIPE, err := this.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	return []arachne.ProtectedEntityInfo{fsIPE}, nil
}

/*
 * Snapshot APIs
 */
func (this FSProtectedEntity) Snapshot(ctx context.Context) (*arachne.ProtectedEntitySnapshotID, error) {
	return nil, nil
}

func (this FSProtectedEntity) ListSnapshots(ctx context.Context) ([]arachne.ProtectedEntitySnapshotID, error) {
	return nil, nil
}
func (this FSProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete arachne.ProtectedEntitySnapshotID) (bool, error) {
	return true, nil
}
func (this FSProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID arachne.ProtectedEntitySnapshotID) (*arachne.ProtectedEntityInfo, error) {
	return nil, nil
}

func (this FSProtectedEntity) GetComponents(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	return make([]arachne.ProtectedEntity, 0), nil
}

func (this FSProtectedEntity) GetID() arachne.ProtectedEntityID {
	return this.id
}

func NewIDFromString(idStr string) vim.ID {
	return vim.ID{
		Id: idStr,
	}
}

func NewVimIDFromPEID(peid arachne.ProtectedEntityID) vim.ID {
	return vim.ID{
		Id: peid.GetID(),
	}
}

func (this FSProtectedEntity) GetDataReader() (io.Reader, error) {
	reader, writer := io.Pipe()
	go runTar(this.root, writer) // Ignore errors until we figure out how to propagate
	return reader, nil

}

func (this FSProtectedEntity) GetMetadataReader() (io.Reader, error) {
	return nil, nil
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
		if (header.Name == "") {
			return nil	// Don't put an empty record for the root
		}
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

func (this *FSProtectedEntityTypeManager) getDataTransports(id arachne.ProtectedEntityID) ([]arachne.DataTransport,
	[]arachne.DataTransport,
	[]arachne.DataTransport, error) {
	dataS3URL := this.s3URLBase + "fs/" + id.String()
	data := []arachne.DataTransport{
		arachne.NewDataTransportForS3(dataS3URL),
	}

	mdS3URL := dataS3URL + ".md"

	md := []arachne.DataTransport{
		arachne.NewDataTransportForS3(mdS3URL),
	}

	combinedS3URL := dataS3URL + ".zip"
	combined := []arachne.DataTransport{
		arachne.NewDataTransportForS3(combinedS3URL),
	}

	return data, md, combined, nil
}
