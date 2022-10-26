package main

import (
	"context"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/aristoteles"
	"github.com/odysseia-greek/plato/aristoteles/configs"
	"github.com/odysseia/periandros/app"
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

func getLeaderConfig(lock *resourcelock.LeaseLock, id string, app *app.PeriandrosHandler) leaderelection.LeaderElectionConfig {
	leaderConfig := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				created, err := app.CreateUser()
				if err != nil {
					glg.Error("error occurred while creating user")
					glg.Error(err)
				}
				if created {
					glg.Infof("created user while being leader: %s", app.Config.SolonCreationRequest.Username)
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
	glg.Info("\n ____   ___  ____   ____   ____  ____   ___    ____   ___   _____\n|    \\ /  _]|    \\ |    | /    ||    \\ |   \\  |    \\ /   \\ / ___/\n|  o  )  [_ |  D  ) |  | |  o  ||  _  ||    \\ |  D  )     (   \\_ \n|   _/    _]|    /  |  | |     ||  |  ||  D  ||    /|  O  |\\__  |\n|  | |   [_ |    \\  |  | |  _  ||  |  ||     ||    \\|     |/  \\ |\n|  | |     ||  .  \\ |  | |  |  ||  |  ||     ||  .  \\     |\\    |\n|__| |_____||__|\\_||____||__|__||__|__||_____||__|\\_|\\___/  \\___|\n                                                                 \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"Περίανδρος δὲ ἦν Κυψέλου παῖς οὗτος ὁ τῷ Θρασυβούλῳ τὸ χρηστήριον μηνύσας· ἐτυράννευε δὲ ὁ Περίανδρος Κορίνθου\"")
	glg.Info("\"Periander, who disclosed the oracle's answer to Thrasybulus, was the son of Cypselus, and sovereign of Corinth\"")
	glg.Info(strings.Repeat("~", 37))

	glg.Debug("creating config")

	baseConfig := configs.PeriandrosConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	periandrosConfig, ok := unparsedConfig.(*configs.PeriandrosConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	var wg sync.WaitGroup
	leaseLockName := fmt.Sprintf("%s-periandros-lock", periandrosConfig.SolonCreationRequest.PodName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	duration := 1 * time.Second
	timeOut := 5 * time.Minute
	handler := app.PeriandrosHandler{Config: periandrosConfig, Duration: duration, Timeout: timeOut}

	lock := handler.Config.Kube.Workload().GetNewLock(leaseLockName, handler.Config.SolonCreationRequest.PodName, handler.Config.Namespace)

	leaderConfig := getLeaderConfig(lock, handler.Config.SolonCreationRequest.PodName, &handler)
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
