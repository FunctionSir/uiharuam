/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-08-27 20:03:36
 * @LastEditTime: 2025-08-28 15:34:16
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /uiharuam/cmdgen.go
 */

package main

import (
	"crypto/sha512"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/FunctionSir/goset"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

//go:embed initdb.sql
var SqlInitDb string

var ErrSomeFileTooLarge error = errors.New("some files too large")

func initSize() int64 {
	return 1024 + ReservedSpace*1024*1024
}

type UAMFile struct {
	ID          string
	Path        string
	Checksum    string
	IsDir       int
	Size        int64
	AddedAt     int64
	TapeBarcode string
	TapeFileNo  int
}

func (f *UAMFile) Values() []any {
	return []any{f.ID, f.Path, f.Checksum, f.IsDir, f.Size, f.AddedAt, f.TapeBarcode, f.TapeFileNo}
}

const SqlInsertToFiles = "INSERT INTO `FILES` VALUES (?,?,?,?,?,?,?,?);"

func writeToDB(f UAMFile) {
	_, err := MetaDBConn.Exec(SqlInsertToFiles, f.Values()...)
	check(err)
}

func cmdGen(cmd *cobra.Command, args []string) {
	if dirExists(WorkingDir) && !isEmptyDir(WorkingDir) {
		fmt.Println("cowardly refusing to do this since working dir is not empty")
		os.Exit(1)
	}
	if !dirExists(WorkingDir) {
		check(os.MkdirAll(WorkingDir, os.ModePerm))
	}
	MetaDB = path.Join(WorkingDir, "meta.db")
	OutDir = path.Join(WorkingDir, "filelists")
	check(os.MkdirAll(OutDir, os.ModePerm))
	if !fileExists(Source) && !dirExists(Source) {
		fmt.Println("Fatal error: source file or dir not found")
		os.Exit(1)
	}
	if dirExists(MetaDB) || fileExists(MetaDB) {
		fmt.Println("cowardly refusing to do this since meta database seems exists")
		os.Exit(1)
	}
	if !isEmptyDir(OutDir) {
		fmt.Println("cowardly refusing to do this since output dir seems not empty")
		os.Exit(1)
	}
	check(openMetaDB())
	defer func() { check(MetaDBConn.Close()) }()
	_, err := MetaDBConn.Exec(SqlInitDb)
	check(err)
	curTape := 1
	curSize := initSize()
	curBC := ""
	curFNo := 0
	var filelist *os.File
	var usedBC = make(goset.Set[string])
	errWalk := filepath.WalkDir(Source, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		size := sizeInTar(d)
		for curBC == "" {
			curBC = promptInputTapeBarcode(curTape)
			curFNo = 0
			if usedBC.Has(curBC) {
				fmt.Println("This tape was already used, use another one instead.")
				curBC = ""
				continue
			}
			filelist, err = os.Create(path.Join(OutDir, fmt.Sprintf("%d_%s_%d.txt", curTape, curBC, curFNo)))
			if err != nil {
				return err
			}
		}
		if size+initSize() > TapeCap*1024*1024 {
			return fmt.Errorf("%w: %s needs about %.3f MiB of space", ErrSomeFileTooLarge, p, float64(size+initSize())/1024/1024)
		}
		for size+curSize > TapeCap*1024*1024 {
			curTape++
			usedBC.Insert(curBC)
			curBC = ""
			for curBC == "" {
				curBC = promptInputTapeBarcode(curTape)
				if usedBC.Has(curBC) {
					fmt.Println("This tape was already used, use another one instead.")
					curBC = ""
				}
			}
			curFNo = 0
			if filelist != nil {
				check(filelist.Close())
			}
			filelist, err = os.Create(path.Join(OutDir, fmt.Sprintf("%d_%s_%d.txt", curTape, curBC, curFNo)))
			if err != nil {
				return err
			}
			curSize = initSize()
		}
		curSize += size
		_, err = filelist.WriteString(p + "\n")
		if err != nil {
			return err
		}
		check(filelist.Sync())
		func() {
			fileID, err := uuid.NewV7()
			check(err)
			if !d.IsDir() {
				hasher := sha512.New()
				f, err := os.Open(p)
				check(err)
				defer func() { _ = f.Close() }()
				_, err = io.Copy(hasher, f)
				check(err)
				stat, err := os.Stat(p)
				check(err)
				writeToDB(UAMFile{
					ID:          fileID.String(),
					Path:        p,
					Checksum:    fmt.Sprintf("%x", hasher.Sum(nil)),
					IsDir:       0,
					Size:        stat.Size(),
					AddedAt:     time.Now().UnixMilli(),
					TapeBarcode: curBC,
					TapeFileNo:  curFNo,
				})
			} else {
				writeToDB(UAMFile{
					ID:          fileID.String(),
					Path:        p,
					Checksum:    "",
					IsDir:       1,
					Size:        0,
					AddedAt:     time.Now().UnixMilli(),
					TapeBarcode: curBC,
					TapeFileNo:  curFNo,
				})
			}
		}()
		fmt.Println(p)
		return nil
	})
	if filelist != nil {
		check(filelist.Close())
	}
	check(errWalk)
	fmt.Println("All done!")
}
