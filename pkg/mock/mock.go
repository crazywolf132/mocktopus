package mock

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/binbandit/mocktopus/pkg/server"
	"github.com/binbandit/mocktopus/pkg/utils"
	"github.com/spf13/viper"
)

type Mock struct {
	mu         *sync.Mutex
	owner      string
	name       string
	Server     *server.Server
	config     *viper.Viper
	mockConfig *viper.Viper
	location   string
	Routes     []RouteConfig
}

func GetMock(name string, config *viper.Viper) (*Mock, error) {
	var owner string = ""
	if strings.Contains(name, "/") {
		utils.Unpack(strings.Split(name, "/"), &owner, &name)
	}
	location, _ := filepath.Abs(
		filepath.Join(
			utils.StripFileFromPath(config.ConfigFileUsed()),
			config.GetString("repos.location"),
			name,
		),
	)
	server, err := server.CreateServer(config)
	if err != nil {
		return nil, err
	}

	return &Mock{
		mu:       &sync.Mutex{},
		owner:    owner,
		name:     name,
		config:   config,
		location: location,
		Server:   server,
	}, nil
}

func hasSuite(name string, config *viper.Viper) (bool, error) {
	fmt.Println(config.GetStringMap("suites"))
	// for val := config.GetStringMap("suites") {
	// 	fmt.Println("Name: ", name)
	// 	fmt.Println("Val: ", val)
	// }
	return true, nil
}
