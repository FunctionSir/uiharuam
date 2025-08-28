/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-08-27 19:38:49
 * @LastEditTime: 2025-08-28 14:45:24
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /uiharuam/main.go
 */

package main

import (
	"database/sql"
	_ "embed"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed desc-cmd-root.txt
var DescCmdRoot string

var WorkingDir string
var OutDir string
var MetaDB string
var TapeCap int64
var ReservedSpace int64
var Source string
var Device string
var MetaDBConn *sql.DB

func main() {
	var rootCmd = &cobra.Command{
		Use:   "uam",
		Short: "Uiharu Archive Manager",
		Long:  strings.TrimSpace(DescCmdRoot),
	}

	var genCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate files list for tar",
		Run:   cmdGen,
	}
	genCmd.Flags().StringVarP(&WorkingDir, "working-dir", "w", "./", "Where stores meta.db and other files")
	genCmd.Flags().StringVarP(&Source, "source", "s", "", "Where files comes")
	genCmd.Flags().Int64VarP(&TapeCap, "tape-cap", "c", 0, "Tape capacity in MiB")
	genCmd.Flags().Int64VarP(&ReservedSpace, "reserved-space", "r", 4, "Reserved space in MiB")
	_ = genCmd.MarkFlagRequired("working-dir")
	_ = genCmd.MarkFlagRequired("tape-cap")

	var writeCmd = &cobra.Command{
		Use:   "write",
		Short: "Write data to tape according to filelists",
		Run:   cmdWrite,
	}
	writeCmd.Flags().StringVarP(&WorkingDir, "working-dir", "w", "./", "Where stores meta.db and other files")
	writeCmd.Flags().StringVarP(&Device, "device", "d", "/dev/nst0", "No rewind tape device for writing tapes")
	_ = writeCmd.MarkFlagRequired("working-dir")

	rootCmd.AddCommand(genCmd, writeCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
