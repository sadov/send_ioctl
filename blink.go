package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

func blink(device string) {
	fd, err := unix.Open(device, unix.O_RDWR|unix.O_DIRECT, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("ioctl")
	unix.Close(fd)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s device count sleep\n", os.Args[0])
		os.Exit(1)
	}

	count, _ := strconv.Atoi(os.Args[2])
	sleep, _ := time.ParseDuration(os.Args[3])

	for i := 1; i <= count; i++ {
		blink(os.Args[1])
		time.Sleep(sleep)
	}
}
