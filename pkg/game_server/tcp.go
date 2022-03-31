package game_server

import (
	"bufio"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/textproto"
	"wheel_of_fortune/pkg/game"
)

type TcpGameServer struct {
	logger      *logrus.Logger
	gameFactory game.GameFactory
}

func NewTcpGameServer(logger *logrus.Logger, factory game.GameFactory) *TcpGameServer {
	return &TcpGameServer{logger: logger, gameFactory: factory}
}

func (s *TcpGameServer) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		s.logger.Warnf("Failed to close connection: %v", err)
	}
}

func (s *TcpGameServer) closeListener(l net.Listener) {
	err := l.Close()
	if err != nil {
		s.logger.Errorf("TCP socket close failed: %v", err)
	}
}

func (s *TcpGameServer) Start(params Params) error {
	l, err := net.Listen("tcp4", net.JoinHostPort("localhost", params.Port))
	if err != nil {
		return wrapServerErr(err)
	}

	defer s.closeListener(l)
	s.logger.Info("Start listening")
	for {
		conn, err := l.Accept()
		if err != nil {
			s.logger.Warnf("Failed to accept connection: %v", err)
			continue
		}
		s.logger.Info("Accept client")
		go s.handleConnection(conn)
	}
}

func (s *TcpGameServer) handleConnection(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	defer s.closeConn(conn)

	gameInstance := s.gameFactory.NewGame()
	err := gameInstance.Start(&TCPMessageSender{writer: writer}, &TCPMessageAcceptor{reader: tp, writer: writer})
	if err != nil && err != io.EOF {
		s.logger.Warnf("Game failed with error: %v", err)
	}
}

type TCPMessageSender struct {
	writer *bufio.Writer
}

func (s *TCPMessageSender) Send(msg string) error {
	writed, err := s.writer.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}
	if writed < len(msg)+1 {
		return errors.New("fail to write all bytes to writer")
	}

	return s.writer.Flush()
}

type TCPMessageAcceptor struct {
	reader *textproto.Reader
	writer *bufio.Writer
}

func (acceptor *TCPMessageAcceptor) Accept() (string, error) {
	_, err := acceptor.writer.WriteString("\n")
	if err != nil {
		return "", err
	}
	err = acceptor.writer.Flush()
	if err != nil {
		return "", err
	}

	return acceptor.reader.ReadLine()
}
