package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dpb587/bosh-compiled-releases/cli/repository"
	"github.com/dpb587/bosh-compiled-releases/cli/server/v1"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

type Serve struct {
	BindHost string `long:"bind-host" description:"Bind host" default:"0.0.0.0"`
	BindPort string `long:"bind-port" description:"Bind port" default:"8080" env:"PORT"`

	Server []string `long:"server" description:"Remote server to query"`
	Local  []string `long:"local" description:"Local path to query"`

	StaticAsset []string `long:"static-asset" description:"Serve a public asset for download"`
}

var _ flags.Commander = Serve{}

func (c Serve) Execute(args []string) error {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout

	repo := repository.NewMultiRepository()

	for _, local := range c.Local {
		files, err := filepath.Glob(local)
		if err != nil {
			log.Fatal("globbing file repositories: ", err)
		}

		for _, file := range files {
			logger.Infof("using local repository: %s", file)

			repo.Attach(repository.NewFileRepository(file))
		}
	}

	for _, server := range c.Server {
		logger.Infof("using remote repository: %s", server)

		repo.Attach(repository.NewServerRepository(server))
	}

	// start

	http.Handle("/v1/resolve", v1.NewResolve(logger, repo))

	for _, file := range c.StaticAsset {
		http.Handle(fmt.Sprintf("/asset/%s", filepath.Base(file)), v1.NewStaticAsset(logger, file))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger := logger.WithFields(logrus.Fields{
			"request.remote_addr": r.RemoteAddr,
			"request.method":      r.Method,
			"request.uri":         r.RequestURI,
		})

		http.Error(w, fmt.Sprintf("404 Not Found"), http.StatusNotFound)
		logger.WithField("response.status", http.StatusNotFound).Debug("not found")
	})

	bind := fmt.Sprintf("%s:%s", c.BindHost, c.BindPort)

	logger.WithFields(logrus.Fields{
		"server.local_addr": bind,
	}).Print("server is ready")

	log.Fatal(http.ListenAndServe(bind, nil))

	return nil
}
