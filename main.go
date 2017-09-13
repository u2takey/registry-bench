package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "registry bench"
	app.Usage = "registry bench"
	app.Action = run
	app.Version = fmt.Sprintf("1.1.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "daemon.mirror",
			Usage:  "docker daemon registry mirror",
			EnvVar: "DOCKER_MIRROR",
		},
		cli.StringFlag{
			Name:   "daemon.workdir",
			Usage:  "docker daemon rworkdir",
			EnvVar: "DOCKER_WORKDIR",
		},
		cli.StringFlag{
			Name:   "daemon.storage-driver",
			Usage:  "docker daemon storage driver",
			EnvVar: "DOCKER_STORAGE_DRIVER",
		},
		cli.StringFlag{
			Name:   "daemon.storage-path",
			Usage:  "docker daemon storage path",
			Value:  "/var/lib/docker",
			EnvVar: "DOCKER_STORAGE_PATH",
		},
		cli.StringFlag{
			Name:   "daemon.bip",
			Usage:  "docker daemon bridge ip address",
			EnvVar: "DOCKER_BIP",
		},
		cli.StringFlag{
			Name:   "daemon.mtu",
			Usage:  "docker daemon custom mtu setting",
			EnvVar: "DOCKER_MTU",
		},
		cli.StringSliceFlag{
			Name:   "daemon.dns",
			Usage:  "docker daemon dns server",
			EnvVar: "DOCKER_DNS",
		},
		cli.BoolFlag{
			Name:   "daemon.insecure",
			Usage:  "docker daemon allows insecure registries",
			EnvVar: "DOCKER_INSECURE",
		},
		cli.BoolFlag{
			Name:   "daemon.ipv6",
			Usage:  "docker daemon IPv6 networking",
			EnvVar: "DOCKER_IPV6",
		},
		cli.BoolFlag{
			Name:   "daemon.debug",
			Usage:  "docker daemon executes in debug mode",
			EnvVar: "DOCKER_DEBUG,DOCKER_LAUNCH_DEBUG",
		},
		cli.BoolFlag{
			Name:   "daemon.off",
			Usage:  "docker daemon executes in debug mode",
			EnvVar: "DOCKER_DAEMON_OFF",
		},
		cli.StringFlag{
			Name:   "dockerfile",
			Usage:  "build dockerfile",
			Value:  "Dockerfile",
			EnvVar: "DOCKER_DOCKERFILE",
		},
		cli.StringFlag{
			Name:   "context",
			Usage:  "build context",
			Value:  ".",
			EnvVar: "DOCKER_CONTEXT",
		},
		cli.StringSliceFlag{
			Name:   "args",
			Usage:  "build args",
			EnvVar: "DOCKER_BUILD_ARGS",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "docker repository",
			EnvVar: "DOCKER_REPO",
		},
		cli.StringFlag{
			Name:   "docker.registry",
			Usage:  "docker registry",
			EnvVar: "DOCKER_REGISTRY",
		},
		cli.StringFlag{
			Name:   "docker.username",
			Usage:  "docker username",
			EnvVar: "DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "docker.password",
			Usage:  "docker password",
			EnvVar: "DOCKER_PASSWORD",
		},
		cli.StringFlag{
			Name:   "docker.email",
			Usage:  "docker email",
			EnvVar: "DOCKER_EMAIL",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
		cli.StringFlag{
			Name:   "httpproxy",
			Usage:  "httpproxy",
			EnvVar: "HTTPPROXY",
			Value:  "",
		},
		cli.IntFlag{
			Name:   "size.start",
			Usage:  "starting image size (in MB)",
			EnvVar: "START_SIZE",
			Value:  10,
		},
		cli.IntFlag{
			Name:   "size.step",
			Usage:  "increase image size each step (in MB)",
			EnvVar: "STEP_SIZE",
			Value:  10,
		},
		cli.IntFlag{
			Name:   "step.count",
			Usage:  "test step count",
			EnvVar: "STEP_COUNT",
			Value:  5,
		},
		cli.IntFlag{
			Name:   "pull.count",
			Usage:  "pull count in each step, if test with cdn/cache, you may set more than 1",
			EnvVar: "PULL_COUNT",
			Value:  1,
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},
		cli.BoolFlag{
			Name:   "randomtag",
			Usage:  "append random tag",
			EnvVar: "RANDOMTAG",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	registry := c.String("docker.registry")
	username := c.String("docker.username")
	password := c.String("docker.password")
	repo := c.String("repo")

	fmt.Println("registry:", registry, " repo:", repo, " username:", username)

	args := c.StringSlice("args")
	if c.String("httpproxy") != "" {
		args = append(args, "HTTP_PROXY="+c.String("httpproxy"))
		args = append(args, "HTTPS_PROXY="+c.String("httpproxy"))
	}

	docker := Docker{
		Login: Login{
			Registry: registry,
			Username: username,
			Password: password,
			Email:    c.String("docker.email"),
		},
		Build: Build{
			Name:       repo,
			Dockerfile: c.String("dockerfile"),
			Context:    c.String("context"),
			Repo:       repo,
			Args:       args,
		},
		Daemon: Daemon{
			Registry:      registry,
			Mirror:        c.String("daemon.mirror"),
			StorageDriver: c.String("daemon.storage-driver"),
			StoragePath:   c.String("daemon.storage-path"),
			Insecure:      c.Bool("daemon.insecure"),
			IPv6:          c.Bool("daemon.ipv6"),
			Debug:         c.Bool("daemon.debug"),
			Bip:           c.String("daemon.bip"),
			DNS:           c.StringSlice("daemon.dns"),
			MTU:           c.String("daemon.mtu"),
			WorkDir:       c.String("daemon.workdir"),
		},
		Debug: c.Bool("debug"),
	}

	count := c.Int("step.count")
	sizestart := c.Int("size.start")
	sizestep := c.Int("size.step")
	pullcount := c.Int("pull.count")

	result := []*TestCase{}
	for i := 0; i < count; i++ {
		size := sizestart + i*sizestep
		t := &TestCase{FileSize: size, PullCount: pullcount}

		err := docker.Exec(t)
		if err != nil {
			fmt.Println("error ", err)
			return err
		}

		result = append(result, t)

		fmt.Printf("case %d done, size %d M, prepare-cost : %f S, pull-cost : %f S [%d pull], push-cost : %f S\n",
			i, t.FileSize, t.Preparecost.Seconds(), t.Pullcost.Seconds(), pullcount, t.Pushcost.Seconds())

	}
	fmt.Println("---------------------------------------------------------------------------------")
	pullsize, pullcost, pushsize, pushcost := 0.0, 0.0, 0.0, 0.0

	for index, t := range result {

		fmt.Printf("case %d, size %d M, prepare-cost : %f S, pull-cost : %f S [%d pull], push-cost : %f S\n",
			index, t.FileSize, t.Preparecost.Seconds(), t.Pullcost.Seconds(), pullcount, t.Pushcost.Seconds())
		pullcost += t.Pullcost.Seconds()
		if index == 0 {
			pullsize += 0
		} else {
			pullsize += float64(result[index-1].FileSize)
		}
		pushcost += t.Pushcost.Seconds()
		pushsize += float64(t.FileSize)
	}
	fmt.Println("---------------------------------------------------------------------------------")
	fmt.Printf("Summary: pull speed %f M/S, push speed %f M/S",
		pullsize*float64(pullcount)/pullcost, pushsize/pushcost)

	return nil
}
