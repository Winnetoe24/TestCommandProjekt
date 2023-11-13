package main

import (
	"Format"
	"bufio"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	exitchnl := make(chan int)
	Bespiele()
	go func() {
		time.Sleep(time.Second * 10)
		println("Exit After Sleep")
		os.Stdin.Close()
		exitchnl <- 0
	}()
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			os.Stderr.WriteString(text + "\n")
			var komm Format.Kommunikation
			json.Unmarshal([]byte(text), &komm)
			switch komm.KommType {
			case Format.PING:
				response := Format.Kommunikation{
					KommType: Format.PING_RESPONSE,
					Data:     make(map[string]string),
				}
				response.Data["PingID"] = komm.Data["PingID"]
				marshal, err := json.Marshal(response)
				if err != nil {
					os.Stderr.WriteString(err.Error())
				} else {
					os.Stdout.Write(marshal)
					os.Stdout.WriteString("\n")
				}
			}
		}
		exitchnl <- 0
	}()
	exitcode := <-exitchnl
	err := os.Stdout.Close()
	if err != nil {
		log.Println(err)
	}
	//println("Close")
	os.Exit(exitcode)
}

func Bespiele() {
	ping := Format.Kommunikation{
		KommType: Format.PING,
		Data:     make(map[string]string),
	}
	ping.Data["PingID"] = strconv.Itoa(rand.Int())
	marshal, _ := json.Marshal(ping)
	os.Stderr.Write(marshal)
	os.Stderr.WriteString("\n")
}
