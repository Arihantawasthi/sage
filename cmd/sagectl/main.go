package main

import (
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/spmp"
)

func main() {
    client := spmp.NewSPMPClient()
    data, err := client.Show()
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while fetching info: %s", err)
    }
    fmt.Println(data)
}
