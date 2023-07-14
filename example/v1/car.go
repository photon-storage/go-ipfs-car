package main

import (
	"bytes"
	"context"
	"fmt"

	car "github.com/photon-storage/go-ipfs-car"
)

func main() {
	b := car.NewBuilder()
	v1car, err := b.Buildv1(
		context.TODO(),
		"./data",
		car.ImportOpts.CIDv1(),
	)
	if err != nil {
		fmt.Printf("Error creating car v1 builder: %v", err)
		return
	}

	v1buf := bytes.Buffer{}
	if err := v1car.Write(&v1buf); err != nil {
		fmt.Printf("Error writing out car v1 format: %v", err)
		return
	}

	fmt.Printf("Car v1 generated, CID =%v, size = %v\n", v1car.Root(), v1buf.Len())
}
