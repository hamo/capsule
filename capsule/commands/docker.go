package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/docker/docker/builder/parser"
)

var (
	dockerBuildCommand = CapsuleCommand{
		Handler:     cmdDockerBuild,
		Description: "Build Docker Image",
	}

	dockerPullCommand = CapsuleCommand{
		Handler:     cmdDockerPull,
		Description: "Pull Docker image",
	}
)

func init() {
	RegisterCommand("docker-build", &dockerBuildCommand)
	RegisterCommand("docker-pull", &dockerPullCommand)
}

func cmdDockerBuild(args []string, cmdEnv *CommandEnv) error {
	dockerBuildFlag := flag.NewFlagSet("docker build command", flag.ExitOnError)

	dockerBuildFlag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: capsule docker-build [options] PATH\n\n")
		dockerBuildFlag.PrintDefaults()
	}

	var (
		flDockerImageName string
	)

	dockerBuildFlag.StringVar(&flDockerImageName, "name", "", "docker image name")

	dockerBuildFlag.Parse(args)

	if dockerBuildFlag.NArg() != 1 {
		dockerBuildFlag.Usage()
		return nil
	}

	dockerFileURI := dockerBuildFlag.Arg(0)
	// FIXME: support URL
	if fi, err := os.Stat(dockerFileURI); err != nil || !fi.Mode().IsRegular() {
		cmdEnv.Logger.Fatalf("can not read file %s", dockerFileURI)
	}

	f, err := os.Open(dockerFileURI)
	if err != nil {
		cmdEnv.Logger.Fatalf("open file %s failed: %s", dockerFileURI, err)
	}

	ast, err := parser.Parse(f)

	cmdEnv.Logger.Fatalf("%+v", ast.Dump())

	return nil

}

func cmdDockerPull(args []string, cmdEnv *CommandEnv) error {
	return nil
}
