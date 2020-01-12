package Plugins

import (
	"../Misc"
	"../Parse"
	"context"
	"net"
	"sync"
	"time"
)

func SubDomain(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup) {
	tch := make(chan int)
	defer func() {
		<-ch
		wg.Done()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(info.Timeout))
	defer cancel()
	go func() {
		ns, err := net.LookupHost(info.Url)
		if err != nil && info.ErrShow {
			Misc.ErrPrinter.Println(err.Error())
		} else if err == nil {
			success += 1
			Misc.SucceedPrinter.Printf("URL:%s HOST:%v", info.Url, ns)
			if info.Output != "" {
				info.OutputTXT()
			}
		}
		tch <- 1
	}()
	select {
	case <-ctx.Done():
		if info.ErrShow {
			Misc.ErrPrinter.Println(info.Url, "Time out")
		}
		return
	case <-tch:
		return
	}
}

func SDConn(info *Misc.HostInfo, ch chan int) {
	stime := time.Now()
	var wg = sync.WaitGroup{}
	prefixs, err := Parse.Readfile(info.UrlFile)
	Misc.CheckErr(err)
	wg.Add(len(prefixs))
	domain := info.Url
	Misc.InfoPrinter.Println("Total length", len(prefixs))
	for _, prefix := range prefixs {
		url := prefix + "." + domain
		info.Url = url
		go SubDomain(*info, ch, &wg)
		ch <- 1
	}
	wg.Wait()
	end := time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:", success)
	Misc.InfoPrinter.Println("Time consumed:", end)
}
