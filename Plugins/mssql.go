package Plugins

import (
	"../Misc"
	"../Parse"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"sync"
	"time"
)

func MsSQL(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {
	config := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;encrypt=disable;timeout=%d",
		info.Host, info.Username, info.Password, info.Port, info.Timeout)
	conn, err := sql.Open("mssql", config)
	if err != nil && info.ErrShow {
		info.PrintFail()
		<-ch
		return
	}
	err = conn.Ping()
	if err != nil && info.ErrShow {
		info.PrintFail()
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

func MssSQLConn(info *Misc.HostInfo, ch chan int) {
	var wg = sync.WaitGroup{}
	var hosts, usernames, passwords []string
	var err error
	stime := time.Now()
	if info.Ports == "" {
		info.Port = MSSQLPORT
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
				go MsSQL(*info, ch, &wg)
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
