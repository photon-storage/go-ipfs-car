# go-ipfs-car

`go-ipfs-car` is a standalone library for packing files into a [Content Addressable aRchives (CAR)](https://ipld.io/specs/transport/car/) file.
The library can be used to generate `.car` file and upload the packaged data into IPFS all at once.

## Features

- Light-weight in-memory operation without bootstraping Kubo core.
- Seamless integration with IPFS network
- Support Car v1 and v2

## Getting Started

### Installation

```sh
go get -u github.com/photon-storage/go-ipfs-car
```

### Example

Packing Car v1 format:

```go
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
```

Packing Car v2 format:
```go
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
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
