package handlers

import (
	"log/slog"
	"net"

	"github.com/AradD7/lightarr/internal/config"
	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/rules"
	"github.com/AradD7/lightarr/internal/wiz"
)

type App struct {
	Cfg      config.Config
	Db       *database.Queries
	Conn     *net.UDPConn
	BulbsMap map[string]*wiz.Bulb
	Rules    []rules.Rule
	Logger   *slog.Logger
}
