package Misc

import (
	"fmt"
	"net"
	"os"
	"time"
)

//Check if port is available
func (h *HostInfo) CheckPort() error {
	InfoPrinter.Println("Checking whether the target port is open")
	addr := fmt.Sprintf("%s:%d", h.Host, h.Port)
	_, err := net.DialTimeout("tcp", addr, (time.Duration(h.Timeout))*time.Second)
	if err != nil {
		WarnPrinter.Println("Port is closed, please confirm whether the target port is open")
		InfoPrinter.Println("If you confirm that the port of this host is open, you can use the -bypass option to bypass port detection")
		return err
	} else {
		InfoPrinter.Println("Port check completed, port open")
		InfoPrinter.Println("Starting")
		return nil
	}
}

//Checking error
func CheckErr(err error) {
	if err != nil {
		ErrPrinter.Println(err.Error())
		os.Exit(0)
	}
}
