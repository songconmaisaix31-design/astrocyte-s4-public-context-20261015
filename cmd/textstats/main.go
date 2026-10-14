package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/songconmaisaix31-design/astrocyte-s4-public-context-20261015/textstats"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(textstats.Count(string(data))); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
