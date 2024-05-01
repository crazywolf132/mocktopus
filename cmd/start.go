package cmd

import (
	"github.com/binbandit/mocktopus/pkg/mock"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <stub>",
	Short: "Starts a mock service by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		config := LoadConfig()
		mock, err := mock.GetMock(name, config)
		if err != nil {
			return err
		}

		err = mock.Update()
		if err != nil {
			return err
		}

		err = mock.LoadConfig()
		if err != nil {
			return err
		}

		err = mock.Discovery()
		if err != nil {
			return err
		}

		err = mock.MakeRoutes()
		if err != nil {
			return err
		}

		mock.Server.Start()

		// mockServer, err := server.CreateServer("./playground/repos/test/config.toml")
		// if err != nil {
		// 	return err
		// }
		// mockServer.AddRoute(server.Route{
		// 	Matchers: []server.RequestMatcher{
		// 		{
		// 			Path:   "/api/hello",
		// 			Method: http.MethodGet,
		// 		},
		// 	},
		// 	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 		w.Header().Set("Content-Type", "application/json")
		// 		w.Write([]byte(`{"message":"hello world"}`))
		// 	}),
		// })
		// mockServer.Start()

		return nil
	},
}
