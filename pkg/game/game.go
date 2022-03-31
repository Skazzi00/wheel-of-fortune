package server

import (
	"bytes"
	"fmt"
	"regexp"
	"wheel_of_fortune/db"
)

type MessageSender interface {
	Send(msg string) error
}

type MessageAcceptor interface {
	Accept() (string, error)
}

type Game interface {
	Start(sender MessageSender, acceptor MessageAcceptor) error
}

type GameFactory interface {
	NewGame() Game
}

type WheelOfFortuneFactory struct {
	wordDb db.WordDB
	tries  int
}

func NewWheelOfFortuneFactory(wordDb db.WordDB, tries int) *WheelOfFortuneFactory {
	return &WheelOfFortuneFactory{wordDb: wordDb, tries: tries}
}

func (w *WheelOfFortuneFactory) NewGame() Game {
	return NewWheelOfFortune(w.wordDb.GetRandomWord(), w.tries)
}

type WheelOfFortune struct {
	word     string
	isOpened []bool
	tries    int
}

func NewWheelOfFortune(word string, initTries int) *WheelOfFortune {
	return &WheelOfFortune{word: word, isOpened: make([]bool, len(word)), tries: initTries}
}

var isLetter = regexp.MustCompile(`^[a-zA-Z]$`).MatchString

func isValidAns(ans string) bool {
	return isLetter(ans)
}

const (
	WelcomeMsg    = "Welcome!\n"
	ClosedChar    = '*'
	InvalidAnswer = "Invalid answer! Your answer must be a letter of the English language. Please try again\n"
	WinMessage    = "Congratulations! You win!\n"
)

func (w *WheelOfFortune) Start(sender MessageSender, acceptor MessageAcceptor) error {
	err := sender.Send(WelcomeMsg)
	if err != nil {
		return err
	}

	for w.tries > 0 {

		gameState := w.getCurrentState()

		err = sender.Send(gameState)
		if err != nil {
			return err
		}

		err = sender.Send(fmt.Sprintf("%v tries left. Enter your letter: ", w.tries))

		answer, err2 := w.getAnswer(sender, acceptor)
		if err2 != nil {
			return err2
		}

		w.updateState(int32(answer[0]))

		if w.isWin() {
			err = sender.Send(WinMessage)
			if err != nil {
				return err
			}
			break
		}

		w.tries--
	}
	return nil
}

func (w *WheelOfFortune) getAnswer(sender MessageSender, acceptor MessageAcceptor) (string, error) {
	for {
		answer, err := acceptor.Accept()
		if err != nil {
			return "", err
		}

		if isValidAns(answer) {
			return answer, nil
		}

		err = sender.Send(InvalidAnswer)
		if err != nil {
			return "", err
		}
	}
}

func (w *WheelOfFortune) getCurrentState() string {
	var msgBuffer bytes.Buffer
	for pos, char := range w.word {
		if w.isOpened[pos] {
			msgBuffer.WriteRune(char)
		} else {
			msgBuffer.WriteRune(ClosedChar)
		}
	}
	msgBuffer.WriteRune('\n')
	return msgBuffer.String()
}

func (w *WheelOfFortune) updateState(openLetter int32) {
	for pos, char := range w.word {
		if char == openLetter {
			w.isOpened[pos] = true
		}
	}
}

func (w *WheelOfFortune) isWin() bool {
	for _, opened := range w.isOpened {
		if !opened {
			return false
		}
	}
	return true
}
