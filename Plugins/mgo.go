package Plugins

import (
	"../Misc"
	"../Parse"
	"fmt"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2"
	"sync"
	"time"
)

func MongoDB(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {

	url := fmt.Sprintf("mongodb://%s:%s@%s:%d", info.Username, info.Password, info.Host, info.Port)
	conn, err := mgo.DialWithTimeout(url, (time.Duration(info.Timeout))*time.Second)
	if err != nil && info.ErrShow {
		info.PrintFail()
		<-ch
		return
	} else if err == nil {
		success += 1
		conn.Close()
		info.PrintSuccess()
		if info.Output != "" {
			info.OutputTXT()
		}
	}
	wg.Done()
	<-ch
}

func MgoConn(info *Misc.HostInfo, ch chan int) {
	var wg = sync.WaitGroup{}
	var hosts, usernames, passwords []string
	var err error
	stime := time.Now()
	if info.Ports == "" {
		info.Port = MONGODBPORT
	} else {
		p, _ := Parse.ParsePort(info.Ports)
		info.Port = p[0]
	}

	hosts, err = Parse.ParseIP(info.Host)
	Misc.CheckErr(err)

	usernames, err = Parse.ParseUser(info)
	Misc.CheckErr(err)

	passwords, err = Parse.ParsePass(info)
	Misc.CheckErr(err)

	wg.Add(len(hosts) * len(usernames) * len(passwords))
	Misc.InfoPrinter.Println("Total length", len(hosts)*len(usernames)*len(passwords))
	for _, host := range hosts {
		for _, user := range usernames {
			for _, pass := range passwords {
				info.Host = host
				info.Username = user
				info.Password = pass
				go MongoDB(*info, ch, &wg)
				ch <- 1
			}
		}
	}
	wg.Wait()
	end := time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:", success)
	Misc.InfoPrinter.Println("Time consumed:", end)
}
