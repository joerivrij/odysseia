package main

import (
	"context"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia/drakon/app"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"os"
	"strings"
	"sync"
	"time"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func getLeaderConfig(lock *resourcelock.LeaseLock, id string, config *app.DrakonHandler) leaderelection.LeaderElectionConfig {
	leaderConfig := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				created, err := config.CreateRoles()
				if err != nil {
					glg.Error("error occurred while creating roles")
					glg.Error(err)
				}
				if created {
					glg.Info("created roles while being leader")
					glg.Info("service exiting after roles have been created")
					os.Exit(0)
				}
			},
			OnStoppedLeading: func() {
				glg.Infof("leader lost: %s", id)
				os.Exit(0)
			},
			OnNewLeader: func(identity string) {
				if identity == id {
					return
				}
				glg.Infof("new leader elected: %s", identity)
			},
		},
	}

	return leaderConfig
}

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PERIANDROS
	glg.Info("\n ___    ____    ____  __  _   ___   ____  \n|   \\  |    \\  /    ||  |/ ] /   \\ |    \\ \n|    \\ |  D  )|  o  ||  ' / |     ||  _  |\n|  D  ||    / |     ||    \\ |  O  ||  |  |\n|     ||    \\ |  _  ||     ||     ||  |  |\n|     ||  .  \\|  |  ||  .  ||     ||  |  |\n|_____||__|\\_||__|__||__|\\_| \\___/ |__|__|\n                                          \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"ἐν τοίνυν τοῖς περὶ τούτων νόμοις ὁ Δράκων φοβερὸν κατασκευάζων καὶ δεινὸν τό τινʼ αὐτόχειρʼ ἄλλον ἄλλου γίγνεσθαι\"")
	glg.Info("\"Now Draco, in this group of laws, marked the terrible wickedness of homicide by banning the offender from the lustral water\"")
	glg.Info(strings.Repeat("~", 37))

	glg.Debug("creating config")

	baseConfig := configs.DrakonConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	drakonConfig, ok := unparsedConfig.(*configs.DrakonConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	var wg sync.WaitGroup
	leaseLockName := fmt.Sprintf("%s-drakon-lock", drakonConfig.PodName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	handler := app.DrakonHandler{Config: drakonConfig}

	lock := handler.Config.Kube.Workload().GetNewLock(leaseLockName, handler.Config.PodName, handler.Config.Namespace)

	leaderConfig := getLeaderConfig(lock, drakonConfig.PodName, &handler)
	leader, err := leaderelection.NewLeaderElector(leaderConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("error electing leader")
	}

	wg.Add(1)
	go func() {
		leader.Run(ctx)
		wg.Done()
	}()

	for {
		glg.Infof("current leader[%s]", leader.GetLeader())
		time.Sleep(5 * time.Second)
	}
}
