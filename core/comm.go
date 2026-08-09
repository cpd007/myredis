package core

import (
	"golang.org/x/sys/unix"
)

type FdComm struct {
	Fd int
}

func (fd FdComm) Read(buffer []byte) (int, error) {
	return unix.Read(fd.Fd, buffer)
}

func (fd FdComm) Write(buffer []byte) (int, error) {
	return unix.Write(fd.Fd, buffer)
}
