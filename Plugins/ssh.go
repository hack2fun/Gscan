package Plugins

import (
	"../Misc"
	"../Parse"
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
)

const SSHPORT = 22

func SSH(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {
	var err error
	config := &ssh.ClientConfig{
		User:            info.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(info.Password)},
		Timeout:         time.Duration(info.Timeout) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ip := fmt.Sprintf("%s:%d", info.Host, info.Port)
	client, err := ssh.Dial("tcp", ip, config)
	if err != nil && info.ErrShow {
		Misc.ErrPrinter.Println(err.Error())
	} else if err == nil {
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

func SSHConn(info *Misc.HostInfo, ch chan int) {
	var hosts, usernames, passwords []string
	var err error
	var wg = sync.WaitGroup{}
	stime := time.Now()
	if info.Ports == "" {
		info.Port = SSHPORT
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
				go SSH(*info, ch, &wg)
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
