package Plugins

import (
	"../Misc"
	"../Parse"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

func MySQL(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&timeout=%ds",
		info.Username, info.Password, info.Host, info.Port, (time.Duration(info.Timeout))*time.Second)
	db, err := sql.Open("mysql", url)
	if err != nil && info.ErrShow {
		info.PrintFail()
	} else {
		err = db.Ping()
		if err != nil && info.ErrShow {
			info.PrintFail()
		} else if err == nil {
			success += 1
			db.Close()
			info.PrintSuccess()
			if info.Output != "" {
				info.OutputTXT()
			}
		}
	}

	wg.Done()
	<-ch
}

func MySQLConn(info *Misc.HostInfo, ch chan int) {
	var hosts, usernames, passwords []string
	var err error
	var wg = sync.WaitGroup{}
	stime := time.Now()
	if info.Ports == "" {
		info.Port = MYSQLPORT
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
				go MySQL(*info, ch, &wg)
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
