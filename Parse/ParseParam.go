package Parse

import (
	"../Misc"
	"flag"
	"fmt"
	"os"
)

var Param = &Misc.HostInfo{}

//Initialize command line arguments
func init() {
	flag.StringVar(&Param.INIFile, "f", "", "configuration file")
	flag.StringVar(&Param.Url, "url", "", "url")
	flag.StringVar(&Param.UrlFile, "uf", "", "Select the path to the url path dictionary")
	flag.BoolVar(&Param.ErrShow, "v", false, "Show details when scanning")
	flag.BoolVar(&Param.Help, "h", false, "Show help")
	flag.StringVar(&Param.Host, "host", "", "IP address of the host you want to scan,for example: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.11,192.168.11.12")
	flag.BoolVar(&Param.Show, "show", false, "Show all scan type")
	flag.StringVar(&Param.Username, "u", "", "Specify a username")
	flag.StringVar(&Param.Password, "p", "", "Specify a password")
	flag.StringVar(&Param.Ports, "port", "", "Select a port,for example: 22 | 1-65535 | 22,80,3306")
	flag.StringVar(&Param.Userfile, "U", "", "Select the path to the username dictionary")
	flag.StringVar(&Param.Passfile, "P", "", "Select the path to the password dictionary")
	flag.StringVar(&Param.Scantype, "m", "", "Select the type you want to scan.If you don't know the scan type and you can add -show to show all scan types")
	flag.IntVar(&Param.Thread, "t", 300, "Set number of threads")
	flag.Int64Var(&Param.Timeout, "w", 2, "Set timeout")
	flag.StringVar(&Param.Output, "o", "", "Save the results of the scan to a file")
	flag.StringVar(&Param.Cookie, "cookie", "", "Set cookie")
	flag.StringVar(&Param.Header, "header", "", "Set http headers (format: JSON)")
	flag.Usage = Myusage
}

func Myusage() {
	fmt.Fprintf(os.Stdout, "Gscan [--host address|--url url] [-p port] [-u username|-U filename] [-uf urlfile] [-p password|-P filename] [-m type] [-t thread] [-w num] [-o output_file] [-v]\n")
	fmt.Fprintf(os.Stdout, "Examples:\n")
	fmt.Fprintf(os.Stdout, "Gscan --host 127.0.0.1 -p 1-65535 -m portscan\n")
	fmt.Fprintf(os.Stdout, "Gscan --host 127.0.0.1 -m ssh -u root -P pass.txt\n")
	fmt.Fprintf(os.Stdout, `Gscan --url http://www.test.com -m urlscan --cookie "PHPSESSID=abc" --header '{"X-FORWARDED-FOR":"test
.com","Referer":"www.baidu.com"}'`)
	fmt.Println()
	fmt.Println("Usage:")
	flag.PrintDefaults()
}
