package game_server

type Params struct {
	Port string
}

type GameServer interface {
	Start(params Params) error
}
