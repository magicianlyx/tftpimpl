package tftpimpl

import "github.com/pkg/errors"

const (
	opRRQ  = uint16(1)
	opWRQ  = uint16(2)
	opDATA = uint16(3)
	opACK  = uint16(4)
	opERR  = uint16(5)
	opOACK = uint16(6)
)

const (
	modeOCTET    = "octet"
	modeMail     = "mail"
	modeNetAscii = "netascii"
)

const (
	/*
	   0         Not defined, see error message (if any).
	   1         File not found.
	   2         Access violation.
	   3         Disk full or allocation exceeded.
	   4         Illegal TFTP operation.
	   5         Unknown transfer ID.
	   6         File already exists.
	   7         No such user.
	   8		 OACK报文错误
	 */
	tftpErrNotDefined        = 0
	tftpErrFileNotFound      = 1
	tftpErrAccessViolation   = 2
	tftpErrDiskFull          = 3
	tftpErrIllegalOperation  = 4
	tftpErrUnknownTransferID = 5
	tftpErrFileAlreadyExists = 6
	tftpErrNoSuchUser        = 7
	tftpErrOACK              = 8
)

const (
	DataBlockSize = 512 // 数据块大小
	DatagramSize  = 516 //报文包大小
)

var (
	ErrParam       = errors.New("args error")            // 参数错误
	ErrDatagram    = errors.New("datagram format error") // 数据报文格式错误
	ErrDataTooLong = errors.New("data block too long")   // 报文数据段过长
)
