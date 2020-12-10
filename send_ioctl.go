package main

import (
	"fmt"
	"github.com/dswarbrick/smart/scsi"
	"os"
)

func main() {
	d, err := scsi.OpenSCSIAutodetect(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer d.Close()
}
