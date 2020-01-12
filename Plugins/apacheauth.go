package Plugins

import (
	"../Misc"
	"../Parse"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

func Apache(info Misc.HostInfo, ch chan int, wg *sync.WaitGroup, client *http.Client) {
	auth := fmt.Sprintf("%s:%s", info.Username, info.Password)
	b64auth := base64.StdEncoding.EncodeToString([]byte(auth))
	res, _ := http.NewRequest("GET", info.Url, nil)
	res.Header.Add("Authorization", "Basic "+b64auth)
	if info.Cookie != "" {
		res.Header.Add("Cookie", info.Cookie)
	}
	if info.Header != "" {
		var header = make(map[string]string)
		err := json.Unmarshal([]byte(info.Header), &header)
		if err != nil {
			Misc.CheckErr(err)
		}
		for k, v := range header {
			res.Header.Add(k, v)
		}
	}
	resp, err := client.Do(res)
	if err != nil && info.ErrShow {
		Misc.ErrPrinter.Println(err.Error())
	} else if err == nil {
		if resp.StatusCode == 200 {
			success += 1
			info.PrintSuccess()
			if info.Output != "" {
				info.OutputTXT()
			}
		} else if resp.StatusCode != 200 && info.ErrShow {
			info.PrintFail()
		}
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	wg.Done()
	<-ch
}

func ApacheConn(info *Misc.HostInfo, ch chan int) {
	var usernames, passwords []string
	var err error
	var wg = sync.WaitGroup{}
	stime := time.Now()

	var client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second * time.Duration(info.Timeout),
			}).DialContext,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if info.UrlFile == "" {
		info.UrlFile = URLFILE
	}

	url, err := Parse.ParseUrl(info.Url)
	Misc.CheckErr(err)
	info.Url = url

	usernames, err = Parse.ParseUser(info)
	Misc.CheckErr(err)

	passwords, err = Parse.ParsePass(info)
	Misc.CheckErr(err)

	wg.Add(len(usernames) * len(passwords))
	Misc.InfoPrinter.Println("Total length", len(usernames)*len(passwords))
	for _, user := range usernames {
		for _, pass := range passwords {
			info.Username = user
			info.Password = pass
			go Apache(*info, ch, &wg, client)
			ch <- 1
		}
	}
	wg.Wait()
	end := time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:", success)
	Misc.InfoPrinter.Println("Time consumed:", end)
}
