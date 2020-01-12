package Plugins

import (
	"../Misc"
	"../Parse"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func MemCached(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {
	realhost := fmt.Sprintf("%s:%d", info.Host, info.Port)
	client, err := net.DialTimeout("tcp", realhost, time.Duration(info.Timeout)*time.Second)
	if err != nil && info.ErrShow {
		Misc.ErrPrinter.Println(err.Error())
		<-ch
		wg.Done()
		return
	} else if err == nil {
		client.SetDeadline(time.Now().Add(time.Duration(info.Timeout)))
		client.Write([]byte("stats\n")) //Set the key randomly to prevent the key on the server from being overwritten
		rev := make([]byte, 1024)
		n, err := client.Read(rev)
		if err != nil && info.ErrShow {
			Misc.ErrPrinter.Println(err.Error())
		} else if strings.Contains(string(rev[:n]), "STAT") {
			success += 1
			client.Close()
			info.PrintSuccess()

		} else {
			Misc.ErrPrinter.Println(string(rev[:n]))
		}
	}
	<-ch
	wg.Done()
	return
}

func MemConn(info *Misc.HostInfo, ch chan int) {
	var hosts []string
	var err error
	var wg = sync.WaitGroup{}
	stime := time.Now()
	if info.Ports == "" {
		info.Port = MEMCACHED
	} else {
		p, _ := Parse.ParsePort(info.Ports)
		info.Port = p[0]
	}

	hosts, err = Parse.ParseIP(info.Host)
	Misc.CheckErr(err)
	wg.Add(len(hosts))
	Misc.InfoPrinter.Println("Total length", len(hosts))
	for _, host := range hosts {
		info.Host = host
		go MemCached(*info, ch, &wg)
		ch <- 1
	}
	wg.Wait()
	end := time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:", success)
	Misc.InfoPrinter.Println("Time consumed:", end)
}
