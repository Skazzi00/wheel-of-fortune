package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"wheel_of_fortune/pkg/db"
	"wheel_of_fortune/pkg/game"
	"wheel_of_fortune/pkg/game_server"
)

const (
	DefaultPort          = "8888"
	DefaultWordsFilename = "/text.txt"
	DefaultTries         = 30
)

type Flags struct {
	WordsFilename *string
	Port          *string
	Tries         *int
}

func getFlags() *Flags {
	flags := &Flags{
		WordsFilename: flag.String("words", DefaultWordsFilename,
			"Filename with word for wheel of fortune, separated by newline"),
		Port:  flag.String("port", DefaultPort, "Port for game_server listen"),
		Tries: flag.Int("tries", DefaultTries, "Number of tries for wheel of fortune game"),
	}
	flag.Parse()
	return flags
}

func getServerParams(flags *Flags) game_server.Params {
	return game_server.Params{Port: *flags.Port}
}

func getDbParams(flags *Flags) db.Params {
	return db.Params{Filename: *flags.WordsFilename, Sep: "\n"}
}

func main() {

	flags := getFlags()
	serverParams := getServerParams(flags)
	dbParams := getDbParams(flags)

	wordsDb, err := db.NewWordDBFromFile(dbParams)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}

	serv := game_server.NewTcpGameServer(logrus.New(), game.NewWheelOfFortuneFactory(wordsDb, *flags.Tries))
	err = serv.Start(serverParams)
	if err != nil {
		panic(err)
		return
	}
}
