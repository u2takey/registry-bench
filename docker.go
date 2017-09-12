package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type (
	// Daemon defines Docker daemon parameters.
	Daemon struct {
		Registry      string   // Docker registry
		Mirror        string   // Docker registry mirror
		Insecure      bool     // Docker daemon enable insecure registries
		StorageDriver string   // Docker daemon storage driver
		StoragePath   string   // Docker daemon storage path
		Running       bool     // Docker daemon is disabled (already running)
		Login         bool     // Docker daemon is loged in
		Debug         bool     // Docker daemon started in debug mode
		Bip           string   // Docker daemon network bridge IP address
		DNS           []string // Docker daemon dns server
		MTU           string   // Docker daemon mtu setting
		IPv6          bool     // Docker daemon IPv6 networking
		WorkDir       string
	}

	// Login defines Docker login parameters.
	Login struct {
		Registry string // Docker registry address
		Username string // Docker registry username
		Password string // Docker registry password
		Email    string // Docker registry email
	}

	// Build defines Docker build parameters.
	Build struct {
		Name       string   // Docker build using default named tag
		Dockerfile string   // Docker build Dockerfile
		Context    string   // Docker build context
		Args       []string // Docker build args
		Repo       string   // Docker build repository
	}

	// Docker defines the Docker plugin parameters.
	Docker struct {
		Login  Login  // Docker login configuration
		Build  Build  // Docker build configuration
		Daemon Daemon // Docker daemon configuration
	}

	// TestCase ...
	TestCase struct {
		FileSize    int
		Preparecost time.Duration
		Pullcost    time.Duration
		Pushcost    time.Duration
	}
)

// Exec executes the plugin step
func (p *Docker) Exec(testcase *TestCase) error {
	start := time.Now()

	// 1. start daemon
	if !p.Daemon.Running {
		cmd := commandDaemon(p.Daemon)
		if p.Daemon.Debug {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			cmd.Stdout = ioutil.Discard
			cmd.Stderr = ioutil.Discard
		}

		go func() {
			trace(cmd)
			cmd.Run()
		}()

		// poll the docker daemon until it is started. This ensures the daemon is
		// ready to accept connections before we proceed.
		for i := 0; i < 15; i++ {
			cmd := commandInfo()
			err := cmd.Run()
			if err == nil {
				p.Daemon.Running = true
				break
			}
			time.Sleep(time.Second * 1)
		}
	}

	start = time.Now()

	// 2. login to the Docker registry
	if p.Daemon.Login {
	} else {
		if p.Login.Password != "" {
			cmd := commandLogin(p.Login)
			trace(cmd)
			var err error
			for i := 0; i < 5; i++ {
				err = cmd.Run()
				if err == nil {
					break
				} else {
					time.Sleep(time.Duration(i+1) * time.Second)
				}
			}
			if err != nil {
				return fmt.Errorf("Error authenticating %s", err)
			}
		} else {
			fmt.Println("Registry credentials not provided. Guest mode enabled.")
		}
		fmt.Println("Registry Login Cost", time.Now().Sub(start))
		p.Daemon.Login = true

		version := commandVersion() // docker version
		info := commandInfo()       // docker info

		version.Run()
		info.Run()
	}

	if p.Daemon.WorkDir != "" {
		err := os.Chdir(p.Daemon.WorkDir)
		if err != nil {
			fmt.Println(err)
		}
	}

	// 3. prepare
	dd := commandDD(testcase.FileSize)
	trace(dd)
	err := dd.Run()
	if err != nil {
		return err
	}

	// 4. build
	start = time.Now()
	target := fmt.Sprintf("%s/%s", p.Daemon.Registry, p.Build.Repo)
	p.Build.Name = target
	build := commandBuild(p.Build)
	trace(build)
	err = build.Run()
	if err != nil {
		return err
	}
	//testcase.Buildcost = time.Now().Sub(start)

	// 5. tag
	// tagcmd, tagname := commandTag(p.Daemon.Registry, p.Build) // docker tag
	// tagcmd.Run()
	// trace(tagcmd)
	// if err != nil {
	// 	return err
	// }

	// time.Sleep(999 * time.Second)

	// 6. push
	start = time.Now()
	push := commandPush(target)
	trace(push)
	err = push.Run()
	if err != nil {
		return err
	}
	testcase.Pushcost = time.Now().Sub(start)

	// 7. rm & pull
	rm := commandRm(target)
	trace(rm)
	err = rm.Run()
	if err != nil {
		return err
	}
	start = time.Now()
	pull := commandPull(target)
	trace(pull)
	err = pull.Run()
	if err != nil {
		return err
	}
	testcase.Pullcost = time.Now().Sub(start)

	return nil
}

const dockerExe = "/usr/local/bin/docker"

// helper function to generate random file
func commandDD(size int) *exec.Cmd {
	return exec.Command(
		"dd", "if=/dev/urandom",
		"of=random.file", "bs="+strconv.Itoa(size)+"M",
		"count=1",
	)
}

// helper function to create the docker login command.
func commandLogin(login Login) *exec.Cmd {
	if login.Email != "" {
		return commandLoginEmail(login)
	}
	return exec.Command(
		dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		login.Registry,
	)
}

func commandLoginEmail(login Login) *exec.Cmd {
	return exec.Command(
		dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		"-e", login.Email,
		login.Registry,
	)
}

// helper function to create the docker info command.
func commandVersion() *exec.Cmd {
	return exec.Command(dockerExe, "version")
}

// helper function to create the docker info command.
func commandInfo() *exec.Cmd {
	return exec.Command(dockerExe, "info")
}

// helper function to create the docker build command.
func commandBuild(build Build) *exec.Cmd {
	cmd := exec.Command(
		dockerExe, "build",
		"--pull=true",
		"--rm=true",
		"--no-cache=true",
		"-f", build.Dockerfile,
		"-t", build.Name,
	)
	for _, arg := range build.Args {
		cmd.Args = append(cmd.Args, "--build-arg", arg)
	}
	cmd.Args = append(cmd.Args, build.Context)
	return cmd
}

// helper function to create the docker tag command.
func commandTag(registry string, build Build) (*exec.Cmd, string) {
	var (
		source = build.Name
		target = fmt.Sprintf("%s/%s", registry, build.Repo)
	)
	return exec.Command(
		dockerExe, "tag", source, target,
	), target
}

// helper function to create the docker push command.
func commandPush(target string) *exec.Cmd {
	return exec.Command(dockerExe, "push", target)
}

func commandRm(target string) *exec.Cmd {
	return exec.Command(dockerExe, "rmi", target)
}

func commandPull(target string) *exec.Cmd {
	return exec.Command(dockerExe, "pull", target)
}

// helper function to create the docker daemon command.
func commandDaemon(daemon Daemon) *exec.Cmd {
	args := []string{"daemon", "-g", daemon.StoragePath}

	if daemon.StorageDriver != "" {
		args = append(args, "-s", daemon.StorageDriver)
	}
	if daemon.Insecure && daemon.Registry != "" {
		args = append(args, "--insecure-registry", daemon.Registry)
	}
	if daemon.IPv6 {
		args = append(args, "--ipv6")
	}
	if len(daemon.Mirror) != 0 {
		args = append(args, "--registry-mirror", daemon.Mirror)
	}
	if len(daemon.Bip) != 0 {
		args = append(args, "--bip", daemon.Bip)
	}
	for _, dns := range daemon.DNS {
		args = append(args, "--dns", dns)
	}
	if len(daemon.MTU) != 0 {
		args = append(args, "--mtu", daemon.MTU)
	}
	return exec.Command(dockerExe, args...)
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
