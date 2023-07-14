package main

import (
	"bytes"
	"context"
	"fmt"

	car "github.com/photon-storage/go-ipfs-car"
)

func main() {
	b := car.NewBuilder()

	ch := make(chan *car.ImportEvent, 8)
	go func() {
		for v := range ch {
			fmt.Printf("New block built, CID = %v\n", v.CID)
		}
		fmt.Printf("Done building\n")
	}()

	v2car, err := b.Buildv2(
		context.TODO(),
		"./data",
		car.ImportOpts.CIDv1(),
		car.ImportOpts.Events(ch),
	)
	if err != nil {
		fmt.Printf("Error creating car v2 builder: %v", err)
		return
	}

	v2buf := bytes.Buffer{}
	if _, err := v2car.WriteTo(&v2buf); err != nil {
		fmt.Printf("Error writing out car v2 format: %v", err)
		return
	}

	fmt.Printf("Car v2 generated, CID =%v, size = %v\n", v2car.Root(), v2buf.Len())
}
