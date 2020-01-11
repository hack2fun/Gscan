package Plugins

import "fmt"

//Plugins List
var PluginList = map[string]interface{}{
	"ftp": FtpConn,
	"mysql": MySQLConn,
	"mongodb":MgoConn,
	"mssql":MssSQLConn,
	"redis": RedisConn,
	"smb": SmbConn,
	"ssh": SSHConn,
	"portscan": PortConn,
	"icmp": IcmpConn,
	"postgresql": PostgreSQLConn,
	"urlscan":UrlConn,
	"auth":ApacheConn,
	"subdomain":SDConn,
	"memcached":MemConn,
}

//Count the number of successes
var success int = 0

const FTPPORT = 21
const MEMCACHED = 11211
const MONGODBPORT  = 27017
const MSSQLPORT = 1433
const PSQLPORT = 5432
const REDISPORT = 6379
const MYSQLPORT  = 3306
const SMBPORT = 445

const URLFILE = "./dict/dicc.txt" //the dict from https://github.com/maurosoria/dirsearch/blob/master/db/dicc.txt
const SUBFILE = "./dict/sub.txt"

//show all Plugins
func Show(){
	fmt.Println("-m")
	for name,_:=range PluginList{
		fmt.Println("   ["+name+"]")
	}
}