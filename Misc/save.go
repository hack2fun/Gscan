package Misc

import (
	"fmt"
	"log"
	"os"
)

//Output the result as text to a file
func (h *HostInfo) OutputTXT() {
	f, err := os.OpenFile(h.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ErrPrinter.Println(err)
		return
	}
	var SaveFiler = log.New(f, "[*]", log.LstdFlags)

	switch h.Scantype {
	case "icmp":
		info := fmt.Sprintf("The host %s is up", h.Host)
		SaveFiler.Println(info)
	case "portscan":
		info := fmt.Sprintf("IP: %s:%d Open", h.Host, h.Port)
		SaveFiler.Println(info)
	//case "mysql","mssql","postgresql","ftp","mongodb","smb","redis","stmp"
	case "urlscan", "subdomain":
		info := fmt.Sprintf("%s", h.Url)
		SaveFiler.Println(info)
	case "auth":
		info := fmt.Sprintf("Type: %s URL: %s Username: %s Password: %s", h.Scantype, h.Url, h.Username, h.Password)
		SaveFiler.Println(info)
	default:
		info := fmt.Sprintf("Type: %s IP: %s:%d Username: %s Password: %s", h.Scantype, h.Host, h.Port, h.Username, h.Password)
		SaveFiler.Println(info)
	}
}
