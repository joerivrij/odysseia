package main

import (
	"context"
	"github.com/kpango/glg"
	"github.com/odysseia/periandros/app"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"os"
	"strings"
	"time"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func getNewLock(lockName, podName, namespace string, client *clientset.Clientset) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockName,
			Namespace: namespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podName,
		},
	}
}

func runLeaderElection(lock *resourcelock.LeaseLock, ctx context.Context, id string, config app.PeriandrosConfig) {
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				glg.Info("creating user")
				created, err := config.CreateUser()
				if err != nil {
					glg.Fatal(err)
				}

				glg.Infof("user created was %t", created)
				return
			},
			OnStoppedLeading: func() {
				glg.Info("no longer the leader, staying inactive.")
			},
			OnNewLeader: func(currentId string) {
				if currentId == id {
					glg.Info("still the leader!")
					return
				}
				glg.Info("new leader is %s", currentId)
			},
		},
	})
}

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PERIANDROS
	glg.Info("\n ____   ___  ____   ____   ____  ____   ___    ____   ___   _____\n|    \\ /  _]|    \\ |    | /    ||    \\ |   \\  |    \\ /   \\ / ___/\n|  o  )  [_ |  D  ) |  | |  o  ||  _  ||    \\ |  D  )     (   \\_ \n|   _/    _]|    /  |  | |     ||  |  ||  D  ||    /|  O  |\\__  |\n|  | |   [_ |    \\  |  | |  _  ||  |  ||     ||    \\|     |/  \\ |\n|  | |     ||  .  \\ |  | |  |  ||  |  ||     ||  .  \\     |\\    |\n|__| |_____||__|\\_||____||__|__||__|__||_____||__|\\_|\\___/  \\___|\n                                                                 \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"Περίανδρος δὲ ἦν Κυψέλου παῖς οὗτος ὁ τῷ Θρασυβούλῳ τὸ χρηστήριον μηνύσας· ἐτυράννευε δὲ ὁ Περίανδρος Κορίνθου\"")
	glg.Info("\"Periander, who disclosed the oracle's answer to Thrasybulus, was the son of Cypselus, and sovereign of Corinth\"")
	glg.Info(strings.Repeat("~", 37))

	glg.Debug("creating config")

	config := app.Get()

	healthy := config.CheckSolonHealth(120)
	if !healthy {
		glg.Fatal("death has found me")
	}

	leaseLockName := "locked-lease"
	leaseLockNamespace := "odysseia"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := rest.InClusterConfig()
	client := clientset.NewForConfigOrDie(cfg)
	if err != nil {
		glg.Fatalf("failed to get incluster kube")
	}

	lock := getNewLock(leaseLockName, config.SolonCreationRequest.PodName, leaseLockNamespace, client)
	runLeaderElection(lock, ctx, config.SolonCreationRequest.PodName, *config)

	os.Exit(0)
}
