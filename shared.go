/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-08-27 20:16:33
 * @LastEditTime: 2025-08-28 14:41:06
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /uiharuam/shared.go
 */
package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"math"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err)
		os.Exit(1)
	}
}

func fileExists(name string) bool {
	stat, err := os.Stat(name)
	if os.IsNotExist(err) || stat.IsDir() {
		return false
	}
	return true
}

func dirExists(name string) bool {
	stat, err := os.Stat(name)
	if os.IsNotExist(err) || !stat.IsDir() {
		return false
	}
	return true
}

func openMetaDB() error {
	var err error
	MetaDBConn, err = sql.Open("sqlite", MetaDB)
	if err != nil {
		return err
	}
	return MetaDBConn.Ping()
}

// As Bytes
func sizeInTar(d fs.DirEntry) int64 {
	if d.IsDir() {
		return 512 // A dir have 512B of header
	}
	info, err := d.Info()
	check(err)
	return 512 + int64(math.Ceil(float64(info.Size())/512)*512)
}

func isEmptyDir(p string) bool {
	entries, err := os.ReadDir(p)
	check(err)
	return len(entries) == 0
}

func promptInputTapeBarcode(curNo int) string {
	tapeBC := ""
	for tapeBC == "" {
		if curNo < 0 {
			fmt.Print("Input the barcode of tape (using scanner is recommended, use ! to terminate): ")
		} else {
			fmt.Printf("Input the barcode of tape No. %d (using scanner is recommended, use ! to terminate): ", curNo)
		}
		_, _ = fmt.Scanln(&tapeBC)
	}
	if tapeBC == "!" {
		os.Exit(0)
	}
	return strings.TrimSpace(tapeBC)
}
