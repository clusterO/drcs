package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// zip a directory
func RecursiveZip(pathToZip, destinationPath string) error {
    destinationFile, err := os.Create(destinationPath)
    if err != nil {
        return err
    }
    myZip := zip.NewWriter(destinationFile)
    err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
        if info.IsDir() {
            return nil
        }
        if err != nil {
            return err
        }
        relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
        zipFile, err := myZip.Create(relPath)
        if err != nil {
            return err
        }
        fsFile, err := os.Open(filePath)
        if err != nil {
            return err
        }
        _, err = io.Copy(zipFile, fsFile)
        if err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return err
    }
    err = myZip.Close()
    if err != nil {
        return err
    }
    return nil
}

// Unzip file to specific destination
func Unzip(src string, dest string) ([]string, error) {

    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {
        fpath := filepath.Join(dest, f.Name)

        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return filenames, fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return filenames, err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return filenames, err
        }

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }

        _, err = io.Copy(outFile, rc)

        outFile.Close()
        rc.Close()

        if err != nil {
            return filenames, err
        }
    }
    
    return filenames, nil
}