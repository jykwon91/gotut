package main

import (
	"flag"
	"log"
	"io/ioutil"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"os"
	"fmt"
	pb "../api"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of the say backend")
	output := flag.String("o", "output.wav", "wav file where")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \" text to speak \"", os.Args[0])
		os.Exit(1)
	}

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to %s: %v", *backend, err)
	}
	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)
	text := &pb.Text{Text: flag.Arg(0)}
	res, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatalf("could not say %s: %v", text.Text, err)
	}
	if err := ioutil.WriteFile(*output, res.Audio, 0666); err != nil {
		log.Fatalf("Could not write to %s: %v", *output, err)
	}

}
