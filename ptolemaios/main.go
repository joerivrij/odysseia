package main

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia/ptolemaios/app"
	"google.golang.org/grpc"
	"net"
	"os"

	pb "github.com/odysseia-greek/plato/proto"
)

const standardPort = ":50051"

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PTOLEMAIOS
	glg.Info("\n ____  ______   ___   _        ___  ___ ___   ____  ____  ___   _____\n|    \\|      | /   \\ | |      /  _]|   |   | /    ||    |/   \\ / ___/\n|  o  )      ||     || |     /  [_ | _   _ ||  o  | |  ||     (   \\_ \n|   _/|_|  |_||  O  || |___ |    _]|  \\_/  ||     | |  ||  O  |\\__  |\n|  |    |  |  |     ||     ||   [_ |   |   ||  _  | |  ||     |/  \\ |\n|  |    |  |  |     ||     ||     ||   |   ||  |  | |  ||     |\\    |\n|__|    |__|   \\___/ |_____||_____||___|___||__|__||____|\\___/  \\___|\n                                                                     \n")
	glg.Info("\"Σωτήρ\"")
	glg.Info("\"savior\"")
	glg.Info("starting up.....")
	glg.Debug("starting up and getting env variables")

	baseConfig := configs.PtolemaiosConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	ptolemaiosConfig, ok := unparsedConfig.(*configs.PtolemaiosConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	handler := app.CreateHandler(*ptolemaiosConfig)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s", port))
	if err != nil {
		glg.Fatalf("failed to listen: %v", err)
	}

	glg.Infof("%s : %s", "setting up rpc service on", port)

	s := grpc.NewServer()
	pb.RegisterPtolemaiosServer(s, handler)
	glg.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		glg.Fatalf("failed to serve: %v", err)
	}
}
