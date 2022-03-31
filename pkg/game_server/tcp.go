package server

import (
	"bufio"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
	"net/textproto"
)

type TcpGameServer struct {
	logger      *logrus.Logger
	gameFactory GameFactory
}

func NewTcpGameServer(logger *logrus.Logger, factory GameFactory) *TcpGameServer {
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

	for {
		conn, err := l.Accept()
		if err != nil {
			s.logger.Warnf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *TcpGameServer) handleConnection(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	defer s.closeConn(conn)

	game := s.gameFactory.NewGame()
	err := game.Start(&TCPMessageSender{writer: writer}, &TCPMessageAcceptor{reader: tp})
	if err != nil {
		s.logger.Warnf("Game failed with error: %v", err)
	}
}

type TCPMessageSender struct {
	writer *bufio.Writer
}

func (s *TCPMessageSender) Send(msg string) error {
	writed, err := s.writer.Write([]byte(msg))
	if err != nil {
		return err
	}
	if writed < len(msg) {
		return errors.New("fail to write all bytes to writer")
	}

	return s.writer.Flush()
}

type TCPMessageAcceptor struct {
	reader *textproto.Reader
}

func (acceptor *TCPMessageAcceptor) Accept() (string, error) {
	return acceptor.reader.ReadLine()
}
