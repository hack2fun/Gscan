package Plugins

import (
	"../Misc"
	"../Parse"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

func Redis(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {
	ip := fmt.Sprintf("%s:%d", info.Host, info.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: info.Password,
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if err != nil && !strings.Contains(pong, "PONG") && info.ErrShow {
		info.PrintFail()
	} else if err == nil && strings.Contains(pong, "PONG") {
		success += 1
		client.Close()
		info.PrintSuccess()
		if info.Output != "" {
			info.OutputTXT()
		}
	}
	wg.Done()
	<-ch
}

func RedisConn(info *Misc.HostInfo, ch chan int) {
	var hosts, passwords []string
	var err error
	var wg = sync.WaitGroup{}
	stime := time.Now()
	if info.Ports == "" {
		info.Port = REDISPORT
	} else {
		p, _ := Parse.ParsePort(info.Ports)
		info.Port = p[0]
	}

	hosts, err = Parse.ParseIP(info.Host)
	Misc.CheckErr(err)

	passwords, err = Parse.ParsePass(info)
	Misc.CheckErr(err)
	wg.Add(len(hosts) * len(passwords))
	Misc.InfoPrinter.Println("Total length", len(hosts)*len(passwords))
	for _, host := range hosts {
		for _, pass := range passwords {
			info.Host = host
			info.Password = pass
			go Redis(*info, ch, &wg)
			ch <- 1
		}
	}
	wg.Wait()
	end := time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:", success)
	Misc.InfoPrinter.Println("Time consumed:", end)
}
