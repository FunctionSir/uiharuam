package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/FunctionSir/ltouwrap"
	"github.com/spf13/cobra"
)

func cmdWrite(cmd *cobra.Command, args []string) {
	if os.Geteuid() != 0 {
		fmt.Println("This command needs root privileges to run.")
		os.Exit(1)
	}
	writtenFilelists := path.Join(WorkingDir, "WRITTEN_FILELISTS")
	var f *os.File
	var err error
	lines := make([]string, 0)
	if !fileExists(writtenFilelists) {
		f, err = os.Create(writtenFilelists)
		check(err)
	} else {
		var tmp []byte
		tmp, err = os.ReadFile(writtenFilelists)
		check(err)
		lines = strings.Split(string(tmp), "\n")
		f, err = os.OpenFile(writtenFilelists, os.O_WRONLY|os.O_APPEND, os.ModePerm)
		check(err)
	}
	defer func() { _ = f.Close() }()
	filelistsDir := path.Join(WorkingDir, "filelists")
	flE, err := os.ReadDir(filelistsDir)
	check(err)
	written := make(map[string]bool)
	for _, x := range flE {
		written[x.Name()] = false
	}
	for _, x := range lines {
		if trimmed := strings.TrimSpace(x); trimmed != "" {
			written[trimmed] = true
		}
	}
	nst, err := ltouwrap.NewLtoNoRewindTapeDrive(Device)
	check(err)
	if _, err := nst.HasDataCartridge(); err == nil {
		fmt.Println("Ejecting current tape...")
		check(nst.Eject(0))
	}
	sort.Slice(flE, func(i, j int) bool {
		iX, err := strconv.Atoi(strings.Split(flE[i].Name(), "_")[0])
		check(err)
		jX, err := strconv.Atoi(strings.Split(flE[j].Name(), "_")[0])
		check(err)
		return iX < jX
	})
	for _, x := range flE {
		if written[x.Name()] {
			continue
		}
		fmt.Printf("Processing filelist %s...\n", x.Name())
		fullpath := path.Join(filelistsDir, x.Name())
		tapeBC := strings.Split(x.Name(), "_")[1]
		tapeFNoS := strings.Split(strings.Split(x.Name(), "_")[2], ".")[0]
		tapeFNoI, err := strconv.Atoi(tapeFNoS)
		check(err)
		fmt.Println("Please find tape with barcode", tapeBC, "then input its barcode below, then insert it.")
		for {
			curBC := promptInputTapeBarcode(-1)
			if curBC == tapeBC {
				break
			}
			fmt.Printf("Barcode mismatch, want %s, got %s. Please try again.\n", tapeBC, curBC)
			if _, err := nst.HasDataCartridge(); err == nil {
				fmt.Println("Ejecting current tape...")
				check(nst.Eject(0))
			}
		}
		fmt.Print("After inserted the tape and get ready, press Enter...")
		_, _ = fmt.Scanln()
		check(nst.Rewind(0))
		for i := 0; i < tapeFNoI; i++ {
			check(nst.NextFile(0))
		}
		cmd := exec.Command("tar", "--no-recursion", "-cvf", nst.DeviceFile, "-T", fullpath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		check(err)
		_, err = f.WriteString(x.Name() + "\n")
		check(err)
		check(f.Sync())
		check(nst.Rewind(0))
		check(nst.Eject(0))
		fmt.Printf("Finished writing %s.\n", tapeBC)
	}
	fmt.Println("All done!")
}
