package main

import (
	"context"
	"github.com/kpango/glg"
	"github.com/odysseia/periandros/app"
	"github.com/odysseia/plato/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"os"
	"strings"
	"sync"
	"time"
)

const testingEnv = "TEST"

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func getLeaderConfig(lock *resourcelock.LeaseLock, id string, config *app.PeriandrosConfig) leaderelection.LeaderElectionConfig {
	leaderConfig := leaderelection.LeaderElectionConfig{
		Lock: lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				created, err := config.CreateUser()
				if err != nil {
					glg.Error("error occurred while creating user")
					glg.Error(err)
				}
				if created {
					glg.Infof("created user while being leader: %s", config.SolonCreationRequest.Username)
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

	env := os.Getenv("ENV")
	if env == "" {
		env = testingEnv
	}

	glg.Infof("env set to: %s", env)

	config := app.Get()

	healthy := config.CheckSolonHealth(120)
	if !healthy {
		glg.Fatal("death has found me")
	}

	var wg sync.WaitGroup
	leaseLockName := "periandros-lock"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	kube := kubernetes.CreateEnvBasedKube(env)
	lock := kube.GetNewLock(leaseLockName, config.SolonCreationRequest.PodName, config.Namespace)

	leaderConfig := getLeaderConfig(lock, config.SolonCreationRequest.PodName, config)
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
