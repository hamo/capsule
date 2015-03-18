package commands

import ()

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
	return nil
}

func cmdDockerPull(args []string, cmdEnv *CommandEnv) error {
	return nil
}
