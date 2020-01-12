package Plugins

import (
	"../Misc"
	"../Parse"
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const count = 1

var reg = regexp.MustCompile(`ttl=\d+`)

func Ping(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {
	var cmd = &exec.Cmd{}
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", strconv.Itoa(count), "-w", strconv.Itoa(int(info.Timeout)), info.Host)
	default:
		cmd = exec.Command("ping", "-c", strconv.Itoa(count), "-W", strconv.Itoa(int(info.Timeout)), info.Host)
	}
	c, err := cmd.StdoutPipe()
	Misc.CheckErr(err)
	defer c.Close()
	cmd.Start()
	result, err := ioutil.ReadAll(c)
	Misc.CheckErr(err)
	tolower := strings.ToLower(string(result))
	ok := reg.MatchString(tolower)
	if ok {
		success += 1
		info.PrintSucceedHost()
		if info.Output != "" {
			info.OutputTXT()
		}
	} else if !ok && info.ErrShow {
		info.PrintFailedHost()
	}
	wg.Done()
	<-ch
}

func IcmpConn(info *Misc.HostInfo, ch chan int) {
	var wg = sync.WaitGroup{}
	stime := time.Now()
	hosts, err := Parse.ParseIP(info.Host)
	Misc.InfoPrinter.Println("Target IP:", info.Host)
	Misc.CheckErr(err)
	wg.Add(len(hosts))
	for _, host := range hosts {
		info.Host = host
		go Ping(*info, ch, &wg)
		ch <- 1
	}
	wg.Wait()
	end := time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:", success)
	Misc.InfoPrinter.Println("Time consumed:", end)
}
