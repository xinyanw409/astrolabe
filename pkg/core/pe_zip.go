package core

import (
	"archive/zip"
	"context"
	"io"
)

func ZipProtectedEntity(ctx context.Context, entity ProtectedEntity, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)
	peInfo, err := entity.GetInfo(ctx)
	if (err != nil) {
		return err
	}
	jsonBuf, err := peInfo.MarshalJSON()
	if (err != nil) {
		return err
	}
	peInfoWriter, err := zipWriter.Create(entity.GetID().String() + ".peinfo")
	if (err != nil) {
		return err
	}
	_, err = peInfoWriter.Write(jsonBuf)
	if (err != nil) {
		return err
	}
	return nil
}