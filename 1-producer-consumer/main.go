//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream, twChan chan *Tweet) {

	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			// stream ended close the channel
			close(twChan)
			return
		}

		twChan <- tweet
	}
}

func consumer(twChan chan *Tweet) {
	for t := range twChan {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	twChan := make(chan *Tweet)
	// Producer
	go producer(stream, twChan)

	// Consumer
	// no need to add 'go' infront of it since main is also a go routine
	consumer(twChan)

	fmt.Printf("Process took %s\n", time.Since(start))
}
