package main

import (
	format "Format"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	exitchnl := make(chan int)
	log.SetFlags(log.Flags() | log.Lshortfile)
	//var outb *bytes.Buffer = bytes.NewBuffer(make([]byte, 200))
	go func() {
		pipe, err2 := GetPipe()
		if err2 != nil {
			log.Fatal(err2)
		}
		inputPipe, err := GetPipe()
		if err != nil {
			log.Fatal(err)
		}
		defer inputPipe.Close()
		command := exec.Command("./CLI/CLI.exe")
		command.Stderr = pipe
		command.Stdout = pipe
		command.Stdin = inputPipe

		if err := command.Start(); err != nil {
			panic(err)
		}
		go func() {
			time.Sleep(time.Second * 12)
			inputPipe.Close()
			pipe.Close()
		}()
		println("After Start")
		scanner := bufio.NewScanner(pipe)

		for scanner.Scan() && !pipe.isClosed.Load() {
			text := scanner.Text()
			fmt.Println(text)
			if text == "Close" {
				err2 := pipe.Close()
				if err2 != nil {
					log.Println(err2)
				}
				break
			}
			var komm format.Kommunikation
			err2 := json.Unmarshal([]byte(text), &komm)
			if err2 != nil {
				log.Println(err2)
				continue
			}
			if komm.KommType == format.PING {
				resp := format.Kommunikation{
					KommType: format.PING_RESPONSE,
					Data:     komm.Data,
				}
				marshal, err2 := json.Marshal(resp)
				if err2 != nil {
					log.Println(err2)
					break
				}
				inputPipe.Write(marshal)
				inputPipe.WriteString("\n")
			}
		}
		println("After Scan")
		if err := scanner.Err(); err != nil && !errors.Is(err, os.ErrClosed) {
			panic(err)
		}
		if err := command.Wait(); err != nil && !errors.Is(err, os.ErrClosed) {
			panic(err)
		}
		exitchnl <- 0
	}()

	exitcode := <-exitchnl
	//if err != nil {
	//	log.Fatal(err)
	//}

	os.Exit(exitcode)
}
