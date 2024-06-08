package instructions

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"redislite/app/data"
	"redislite/app/data/datatypes/kvstream"
	// "redislite/app/data/datatypes/kvstream"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

type streamPair struct {
	Name    string
	StartId string
}

func XRead(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	//redis-cli XREAD block 2000 some_key 1526985054069 1526985054079
	var blockTimeMs int64 = -1
	var err error
	var onlyNewest = false

	startStreamIndex := 0
	dataItems := redisCommand.Params
	firstCommandParam := dataItems[0]

	if strings.ToUpper(firstCommandParam) == "BLOCK" {
		startStreamIndex = 2
		blockTimeMs, err = strconv.ParseInt(dataItems[1], 10, 64)
		if err != nil {
			return errors.New("Malformed string in block command")
		}
	}

	//make sure the streams command is in the params list before we start to process data
	if strings.ToUpper(dataItems[startStreamIndex]) != "STREAMS" {
		return errors.New("Missing STREAMS command")
	}

	streams := []streamPair{}
	streams, err = parseStreamParams(dataItems[startStreamIndex+1:])

	//standard process, get it out of the way now
	if blockTimeMs < 0 {
		result, err := parseStreamPairs(streams, server, false)
		if err != nil {
			return err
		}
		log.Printf("%#v", result)

		return protomessages.QuickSendMessage(conn, result)
	}

	streamName := streams[0].Name
	if streams[0].StartId == "$" {
		log.Println("found start Id")
		onlyNewest = true
		//only want new requests
	}

	log.Println(onlyNewest)
	//BLOCK BLOCK BLOCK
	currentTimeNano := time.Now().UTC().UnixNano()
	firstRecieved := streamPair{}
	recievedMessage := false
mainloop:
	for {
		select {
		case item := <-server.StreamChannel:
			if item.AddedTime < currentTimeNano {
				continue
			}
			if strings.ToUpper(streamName) != strings.ToUpper(item.StreamName) {
				continue
			}
			if !recievedMessage {
				firstRecieved = streamPair{
					Name:    streamName,
					StartId: item.Data.Id.ToString(),
				}
				recievedMessage = true
				if blockTimeMs == 0 {
					break mainloop
				}
			}
		case <-time.After(time.Duration(blockTimeMs * int64(time.Millisecond))):
			if blockTimeMs > 0 {
				break mainloop
			}
		}
	}

	result := ""

	if onlyNewest {
		log.Printf("Parsing newest: %#v", firstRecieved)
		result, err = parseStreamPairs([]streamPair{firstRecieved}, server, true)
		return protomessages.QuickSendMessage(conn, result)
	}

	log.Printf("Parsing all: %#v", streams)
	result, err = parseStreamPairs(streams, server, false)

	if err != nil {
		return err
	}
	return protomessages.QuickSendMessage(conn, result)
}

func parseStreamPairs(streamPairs []streamPair, server *setup.Server, inclusive bool) (string, error) {
	streamContainers := []string{}
	for _, stream := range streamPairs {
		//get back all of the data from an individual stream
		log.Printf("%#v", stream)
		streamrange, err := server.DataStore.GetStreamRead(stream.Name, stream.StartId, inclusive)
		if err != nil {
			return "", err
		}

		dataArray := parseStreamRange(streamrange)

		if len(dataArray) == 0 {
			continue
		}
		//build a redis array with each container
		streamcontainer := protomessages.BuildMultiRespArrayMsg(dataArray)
		//add the container to the stream name
		// streamName = protomessages.BuildBulkStringMsg(streamName)
		streamarray := append([]string{protomessages.BuildBulkStringMsg(stream.Name)}, streamcontainer)
		//add all the containers to the streamcontainer
		streamcontainer = protomessages.BuildMultiRespArrayMsg(streamarray)
		//add the stream container to the final output
		streamContainers = append(streamContainers, streamcontainer)
		//this will eventually be an array of containers
	}
	if len(streamContainers) == 0 {
		return "$-1\r\n", nil
	}
	return protomessages.BuildMultiRespArrayMsg(streamContainers), nil
}

func parseStreamRange(streamrange []kvstream.KvStreamData) []string {
	//the stream name
	dataArray := []string{}
	//for each item in the stream (id + timeseries) add it to an array called dataArray
	for _, streamItem := range streamrange {
		title := protomessages.BuildBulkStringMsg(streamItem.Id.ToString())
		kvparray := []string{}
		//add each kv array to the data arary
		for _, kvp := range streamItem.Items {
			kvparray = append(kvparray, kvp.Key)
			kvparray = append(kvparray, kvp.Value)
		}

		//create a container for the kv array
		kvpresp := protomessages.BuildRespArrayMsg(kvparray)
		//create a container with the title and the kv array
		kvcontainer := protomessages.BuildMultiRespArrayMsg([]string{title, kvpresp})
		dataArray = append(dataArray, kvcontainer)

	}

	return dataArray
}

func parseStreamParams(params []string) ([]streamPair, error) {
	paramlen := len(params)
	if paramlen%2 != 0 {
		return []streamPair{}, errors.New("Wrong number of arguments")
	}
	streamPairs := []streamPair{}
	count := paramlen / 2
	for index := 0; index < count; index++ {
		streamPairs = append(streamPairs, streamPair{
			Name:    params[index],
			StartId: params[index+count],
		})
	}

	return streamPairs, nil
}
