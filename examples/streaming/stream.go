package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/oriiolabs/revai-go"
)

func main() {

	rev_token := os.Getenv("REV_TOKEN")
	if rev_token == "" {
		fmt.Println("Please set REV_TOKEN enviroment variable")
		return
	}

	c := revai.NewClient(rev_token)

	// select the content type for your audio
	params := &revai.DialStreamParams{
		ContentType: "audio/x-wav",
	}

	conn, err := c.Stream.Dial(context.Background(), params)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			msg, err := conn.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Reader got EOF")
					break
				} else {
					panic(err)
				}
			}
			res, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(res))
		}
	}()

	f, err := os.Open("test_data.wav")

	if err != nil {
		panic(err)
	}

	// update the match the byte rate of your audio
	byte_rate := int64(64000) // bit rate / 8
	buffer := make([]byte, byte_rate)
	for {
		if _, err := f.Read(buffer); err != nil {
			if err == io.EOF {
				conn.WriteDone()
				time.Sleep(5 * time.Second)
				break
			}
		} else if err := conn.Write(bytes.NewReader(buffer)); err != nil {
			panic(err)
		}

		time.Sleep(900 * time.Millisecond)

	}

	conn.Close()

}
