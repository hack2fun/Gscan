package Misc

import(
	"log"
	"os"
	"fmt"
)

const REDCOLOR = "\x1B[0;40;31m[Failed] \x1B[0m"
const GREEN = "\x1B[1;40;32m[Succeed] \x1B[0m"
const ERR = "\x1B[0;40;31m[ERROR] \x1B[0m"
const INFO = "\x1B[1;40;32m[*] \x1B[0m"
const WARN = "\x1B[0;40;31m[*] \x1B[0m"

var SucceedPrinter = log.New(os.Stdout,GREEN,log.LstdFlags)
var FailedPrinter = log.New(os.Stdout,REDCOLOR,log.LstdFlags)
var ErrPrinter = log.New(os.Stdout,ERR,log.LstdFlags)
var InfoPrinter = log.New(os.Stdout,INFO,log.LstdFlags)
var WarnPrinter = log.New(os.Stdout,WARN,log.LstdFlags)
// Successful IP, account and password will be output here
func (h *HostInfo)PrintSuccess(){
	info:=fmt.Sprintf("Type: %s IP: %s:%d Username: %s Password: %s",h.Scantype,h.Host,h.Port,h.Username,h.Password)
	SucceedPrinter.Println(info)
}

func (h *HostInfo)PrintFail(){
	info:=fmt.Sprintf("Type: %s IP: %s:%d Username: %s Password: %s",h.Scantype,h.Host,h.Port,h.Username,h.Password)
	FailedPrinter.Println(info)
}

//Output alive ports
func (h *HostInfo)PrintSucceedPort(){
	info:=fmt.Sprintf("IP: %s:%d",h.Host,h.Port)
	SucceedPrinter.Println(info)
}

//Output died ports
func (h *HostInfo)PrintFailedPort(){
	info:=fmt.Sprintf("IP: %s:%d",h.Host,h.Port)
	FailedPrinter.Println(info)
}

//Output alive hosts
func (h *HostInfo)PrintSucceedHost(){
	info:=fmt.Sprintf("IP: %s is aliving",h.Host)
	SucceedPrinter.Println(info)
}

//Output died hosts
func (h *HostInfo)PrintFailedHost(){
	info:=fmt.Sprintf("IP: %s is died",h.Host)
	FailedPrinter.Println(info)
}