package main

import (
	"os"
	"sync/atomic"
)

type SelfClosingPipe struct {
	read     *os.File
	write    *os.File
	isClosed atomic.Bool
}

func GetPipe() (*SelfClosingPipe, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	return &SelfClosingPipe{
		read:  r,
		write: w,
	}, nil

}

func (scp *SelfClosingPipe) Read(p []byte) (n int, err error) {
	println("Read SelfClosingPipe")
	return scp.read.Read(p)
}

func (scp *SelfClosingPipe) Write(p []byte) (n int, err error) {
	println("Wrote SelfClosingPipe")
	return scp.write.Write(p)
}

func (scp *SelfClosingPipe) WriteString(s string) (n int, err error) {
	return scp.write.WriteString(s)
}

func (scp *SelfClosingPipe) Close() error {
	if scp.isClosed.Load() {
		println("SelfClosingPipe Already Closed")
		return nil
	}
	scp.isClosed.Store(true)
	println("SelfClosingPipe Close")
	err := scp.read.Close()
	err2 := scp.write.Close()
	if err == nil {
		return err2
	}
	return err
}
