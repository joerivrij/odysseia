package main

import (
	"context"
	"github.com/kpango/glg"
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

func getLeaderConfig(lock *resourcelock.LeaseLock, id string, config *app.DrakonConfig) leaderelection.LeaderElectionConfig {
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

	config := app.Get()

	var wg sync.WaitGroup
	leaseLockName := "drakon-lock"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lock := config.Kube.GetNewLock(leaseLockName, config.Podname, config.Namespace)

	leaderConfig := getLeaderConfig(lock, config.Podname, config)
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
