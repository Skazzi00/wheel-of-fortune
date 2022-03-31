package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"net/textproto"
)

func startGame(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

outer:
	for {
		for {
			line, err := tp.ReadLine()
			if err != nil {
				if err == io.EOF {
					break outer
				}
				return err
			}
			if len(line) == 0 {
				break
			}
			fmt.Println(line)
		}

		var line string
		_, err := fmt.Scanf("%s", &line)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(conn, "%s\n", line)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	serverAddr := flag.String("s", "localhost:8888", "Address of game server")
	isHelp := flag.Bool("help", false, "Show usage")

	flag.Parse()

	if *isHelp {
		flag.Usage()
		return
	}

	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = startGame(conn)
	if err != nil {
		panic(err)
	}
}
