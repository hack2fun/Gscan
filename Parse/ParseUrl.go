package Parse

import (
	"errors"
	"strings"
)

var UrlErr = errors.New("url parse error" +
	"format:\n" +
	"www.baidu.com\n" +
	"https://www.baidu.com\n"+
	"http://www.baidu.com\n" +
	"http://192.168.1.1\n"+
	"192.168.1.1\n")

func ParseUrl(url string)(string,error){
	if url==""{
		return "",errors.New("Please specify a URL with the \"--url\"")
	}
	switch {
	case strings.HasPrefix(url,"http://") ||
		 strings.HasPrefix(url,"https://") :
		 	if url[len(url)-1]=='/'{
				return url,nil
			}else{
				return url+"/",nil
			}
	case !strings.HasPrefix(url,"http://") &&
		!strings.HasPrefix(url,"https://"):
		if url[len(url)-1]=='/' {
			return "http://"+url,nil
		}else{
			return "http://"+url+"/",nil
		}

	default:
		return "",UrlErr
	}
}
