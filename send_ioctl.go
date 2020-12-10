package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/dswarbrick/smart/ioctl"
	"github.com/dswarbrick/smart/scsi"
)

// SCSI generic ioctl header, defined as sg_io_hdr_t in <scsi/sg.h>
type sgIoHdr struct {
	interface_id    int32   // 'S' for SCSI generic (required)
	dxfer_direction int32   // data transfer direction
	cmd_len         uint8   // SCSI command length (<= 16 bytes)
	mx_sb_len       uint8   // max length to write to sbp
	iovec_count     uint16  // 0 implies no scatter gather
	dxfer_len       uint32  // byte count of data transfer
	dxferp          uintptr // points to data transfer memory or scatter gather list
	cmdp            uintptr // points to command to perform
	sbp             uintptr // points to sense_buffer memory
	timeout         uint32  // MAX_UINT -> no timeout (unit: millisec)
	flags           uint32  // 0 -> default, see SG_FLAG...
	pack_id         int32   // unused internally (normally)
	usr_ptr         uintptr // unused internally
	status          uint8   // SCSI status
	masked_status   uint8   // shifted, masked scsi status
	msg_status      uint8   // messaging level data (optional)
	sb_len_wr       uint8   // byte count actually written to sbp
	host_status     uint16  // errors from host adapter
	driver_status   uint16  // errors from software driver
	resid           int32   // dxfer_len - actual_transferred
	duration        uint32  // time taken by cmd (unit: millisec)
	info            uint32  // auxiliary information
}

func prepare_hdr(fd int) (sgIoHdr, error) {
	senseBuf := make([]byte, 32)

	respBuf := make([]byte, scsi.INQ_REPLY_LEN)
	cdb := scsi.CDB6{scsi.SCSI_INQUIRY}

	binary.BigEndian.PutUint16(cdb[3:], uint16(len(respBuf)))

	// Populate required fields of "sg_io_hdr_t" struct
	hdr := sgIoHdr{
		interface_id:    'S',
		dxfer_direction: scsi.SG_DXFER_FROM_DEV,
		timeout:         scsi.DEFAULT_TIMEOUT,
		cmd_len:         uint8(len(cdb)),
		mx_sb_len:       uint8(len(senseBuf)),
		dxfer_len:       uint32(len(respBuf)),
		dxferp:          uintptr(unsafe.Pointer(&(respBuf)[0])),
		cmdp:            uintptr(unsafe.Pointer(&cdb[0])),
		sbp:             uintptr(unsafe.Pointer(&senseBuf[0])),
	}

	return hdr, nil
}

func send_ioctl(fd int, hdr sgIoHdr) error {
	err := ioctl.Ioctl(uintptr(fd), scsi.SG_IO,
		uintptr(unsafe.Pointer(&hdr)))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var hdr sgIoHdr

	fd, err := unix.Open(os.Args[1], unix.O_RDWR, 0600)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if hdr, err = prepare_hdr(fd); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if err = send_ioctl(fd, hdr); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	defer unix.Close(fd)
}
