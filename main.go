package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

var (
	folder_name  string
	archive_name string
)

type FromTo struct {
	From string
	To   string
}

func main() {
	flag.StringVar(&folder_name, "dir", "", "Folder to zip images in")
	flag.StringVar(&archive_name, "out", "", "Name of zip directory")
	flag.Parse()

	if folder_name == "" {
		panic(`Directory entry malformed or unspecified. Use -dir "/folder/path/goes/here ."`)
	}
	if archive_name == "" || (strings.Split(archive_name, ".")[1] != "zip") {
		panic(`-out not specified or malformed, use -out name.zip`)
	}

	folder_contents, err := os.ReadDir(folder_name)
	if err != nil {
		panic(err)
	}
	// filter for image files from likely to more unlikely
	allowed_formats := []string{"png", "exr", "jpg", "jpeg", "JPG", "JPEG", "EXR", "tiff", "TIFF"}
	files := make([]FromTo, 0)

	for _, entry := range folder_contents {
		if entry.IsDir() {
			continue
		}
		// split so fucking nonsense cases like EXR being present in the filename won't match.
		// grab last element by splitting and reversing and grabbing first
		fname := entry.Name()
		dotsplt := strings.Split(fname, ".")
		slices.Reverse(dotsplt)
		exit := dotsplt[0]

		if slices.Contains(allowed_formats, exit) {
			files = append(files, FromTo{From: folder_name + "/" + fname, To: fname})
		}
	}

	zipfile, err := os.Create(folder_name + "/" + archive_name)
	if err != nil {
		panic(err)
	}
	defer zipfile.Close()

	zwrite := zip.NewWriter(zipfile)
	defer zwrite.Close()

	for idx, fromto := range files {
		source, err := os.Open(fromto.From)
		defer source.Close()
		if err != nil {
			continue
		}

		destination, err := zwrite.Create(fromto.To)
		if err != nil {
			continue
		}

		_, err = io.Copy(destination, source)
		if err != nil {
			os.Stderr.WriteString("Failed to write " + fromto.To + " in archive...\n")
			continue
		}
		os.Stdout.WriteString(fmt.Sprint(idx+1) + "/" + fmt.Sprint(len(files)) + "\n")
	}
	os.Stdout.WriteString("Done! Written to:\n")
	os.Stdout.WriteString(folder_name + "/" + archive_name + "\n")
}
