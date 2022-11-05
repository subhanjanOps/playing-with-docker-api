package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {

	introFigure := figure.NewColorFigure("Container Ops", "doom", "green", true)
	introFigure.Print()

	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	ctx := context.Background()
	// defer cancel()

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

	LogContainer(ctx, cli, containers[inp])

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
	if selection >= 1 && selection <= max {
		return selection
	} else {
		fmt.Println("Selection is out of scope.")
		return GetSelection(max, exitSignal)
	}
}

func LogContainer(ctx context.Context, client *client.Client, container types.Container) error {

	readCloser, err := client.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       "40",
		Details:    true,
	})
	if err != nil {
		return err
	}

	// _, err = io.Copy(os.Stdout, readCloser)
	// if err != nil && err != io.EOF {
	// 	return err
	// }
	// fmt.Println(written)
	// data, _ := io.ReadAll(readCloser)
	// fmt.Println(data)
	// json.NewDecoder(readCloser).Decode(&logOutput)

	logBytes := make([]byte, 8)
	for {
		_, err := readCloser.Read(logBytes)
		if err != nil {
			return err
		}
		var w io.Writer
		switch logBytes[0] {
		case 1:
			w = os.Stdout
		default:
			w = os.Stderr
		}
		count := binary.BigEndian.Uint32(logBytes[4:])
		dat := make([]byte, count)
		_, err = readCloser.Read(dat)
		if err != nil {
			return err
		}
		fmt.Fprint(w, string(dat))
	}

	return nil
}
