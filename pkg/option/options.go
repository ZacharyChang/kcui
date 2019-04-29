package option

import (
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

// Options contains start options
type Options struct {
	Kubeconfig string
	Namespace  string
	Debug      bool
}

// NewOptions returns a new Options
func NewOptions() *Options {
	return &Options{}
}

// AddFlags add flags to app
func (opts *Options) AddFlags(app *cli.App) {

	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "Path to the kube config. This can be set from environment value 'KUBECONFIG'",
			Destination: &opts.Kubeconfig,
			Value:       filepath.Join(os.Getenv("HOME"), ".kube", "config"),
			EnvVar:      "KUBECONFIG",
		},
		cli.StringFlag{
			Name:        "namespace",
			Usage:       "Namespace to query.",
			Destination: &opts.Namespace,
			Value:       "",
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Run with debug mode",
			Destination: &opts.Debug,
		},
	}

	app.Flags = append(app.Flags, flags...)

}
