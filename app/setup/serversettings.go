package setup

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"redislite/app/data"
	"redislite/app/data/storage"
	"redislite/app/persistence"
)

type ServerSettings struct {
	Host               string
	Port               int
	Master             bool
	MasterReplId       string
	MasterReplIdOffset int
	MasterHost         string
	MasterPort         int
	Dir                string
	DbFilename         string
}

type Server struct {
	Replicas        []net.Conn
	DataStore       storage.DataStore
	Settings        ServerSettings
	RecievedCounter *data.ByteCounter
	Rdb             persistence.RdbFile
	AckChannel      chan bool
  ClientChannel   chan Client
}

func CreateServer() *Server {
	settings := CreateServerSettings()
	datastore := storage.CreateKvStore()
	rdb := persistence.RdbFile{}
	if settings.DbFilename != "" {
		rdb = persistence.ReadRdbFromFile(settings.DbFilename, settings.Dir)
		datastore.LoadFromDb(rdb.DataItems)
	}

	return &Server{
		DataStore:       datastore,
		Settings:        settings,
		Replicas:        []net.Conn{},
		RecievedCounter: &data.ByteCounter{},
		Rdb:             rdb,
		AckChannel:      make(chan bool),
    ClientChannel:   make(chan Client),
	}
}

func CreateServerSettings() ServerSettings {
	serverSettings := ServerSettings{
		Host:               "127.0.0.1:6379",
		Port:               6379,
		Master:             true,
		MasterReplId:       "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplIdOffset: 0,
	}
	replicaOf := ""
	flag.IntVar(&serverSettings.Port, "port", 6379, "assigns the redis port")
	flag.StringVar(&replicaOf, "replicaof", "", "assigns a master to be a replica of")
	flag.StringVar(&serverSettings.Dir, "dir", "", "directory of db")
	flag.StringVar(&serverSettings.DbFilename, "dbfilename", "", "filename of db")
	flag.Parse()

	serverSettings.Host = fmt.Sprintf("127.0.0.1:%d", serverSettings.Port)
	if replicaOf != "" {
		log.Println("launching in REPLICA mode")
		host, port, err := parseParamsForReplica(replicaOf)
		if err != nil {
			errmess := fmt.Sprintf("Could not parse replicaOf parameters, malformed string: %s \r\n", err)
			log.Printf(errmess)
			panic(errmess)
		}
		log.Printf("Acting as replica, attempting to connecting to master located at: %s:%d\r\n", host, port)
		serverSettings.MasterHost = host
		serverSettings.MasterPort = port
		serverSettings.Master = false

	}

	return serverSettings
}

func parseParamsForReplica(replicaOf string) (string, int, error) {
	var parts = strings.Split(replicaOf, " ")
	if len(parts) != 2 {
		return "", 0, errors.New("Wrong number of parameters")
	}

	var host string
	hostip := net.ParseIP(parts[0])
	if hostip != nil {
		host = string(hostip.String())
		log.Printf("host ip set, using: %s \r\n", host)
	} else {
		host = parts[0]
		log.Printf("host is not an ip address, attempting to use hostname: %s \r\n", host)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, errors.New("Port is not a number")
	}

	return host, port, nil

}
