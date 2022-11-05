package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("\nNo. ContainerID Name\n")

	for i, container := range containers {
		fmt.Printf("%v. %s %v\n", i+1, container.ID[:10], container.Names[0][1:])
	}
	fmt.Printf("%v. Exit\n", 99)

	inp := GetSelection(len(containers), 99) - 1
	fmt.Println(fmt.Printf("%v. %s %v\n", inp+1, containers[inp].ID[:10], containers[inp].Names[0][1:]))

}
func GetSelection(max int, exitSignal int) int {
	r := bufio.NewReader(os.Stdin)
	var s string
	for {
		fmt.Fprint(os.Stderr, "Enter selection: ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	selection, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	if selection == exitSignal {
		os.Exit(0)
	}
	if selection < 1 || selection > max || selection != exitSignal {
		fmt.Println("Selection is out of scope.")
		return GetSelection(max, exitSignal)
	}
	return selection
}

func LogContainer(container types.Container) error {
	return nil
}
