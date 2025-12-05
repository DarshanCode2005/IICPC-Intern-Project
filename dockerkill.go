// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// )

// func runCommand(command string, args ...string) {
// 	cmd := exec.Command(command, args...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// }

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: dockerkill <command>")
// 		os.Exit(1)
// 	}

// 	switch os.Args[1] {
// 	case "everything":
// 		runCommand("docker", "system", "prune", "-a", "--volumes", "-f")
// 	case "images":
// 		runCommand("docker", "rmi", "$(docker images -q)", "--force")
// 	case "containers":
// 		runCommand("docker", "kill", "$(docker ps -q)")
// 	case "networks":
// 		runCommand("docker", "network", "prune", "-f")
// 	case "volumes":
// 		runCommand("docker", "volume", "prune", "-f")
// 	case "list":
// 		if len(os.Args) < 3 {
// 			fmt.Println("Usage: dockerkill list <images|containers|networks|volumes>")
// 			os.Exit(1)
// 		}
// 		switch os.Args[2] {
// 		case "images":
// 			runCommand("docker", "images")
// 		case "containers":
// 			runCommand("docker", "ps", "-a")
// 		case "networks":
// 			runCommand("docker", "network", "ls")
// 		case "volumes":
// 			runCommand("docker", "volume", "ls")
// 		default:
// 			fmt.Println("Unknown list command.")
// 			os.Exit(1)
// 		}
// 	case "prune":
// 		runCommand("docker", "system", "prune", "-a")
// 	default:
// 		fmt.Println("Unknown command.")
// 		os.Exit(1)
// 	}
// }

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getDockerIDs(args ...string) []string {
	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	lines := strings.Split(string(out), "\n")
	ids := []string{}
	for _, line := range lines {
		id := strings.TrimSpace(line)
		if len(id) > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}

func runDockerRemove(ids []string, removeCmd string) {
	for _, id := range ids {
		fmt.Println(removeCmd, id)
		cmd := exec.Command("docker", removeCmd, id)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dockerkill <command>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "killcontainers":
		ids := getDockerIDs("ps", "-q")
		runDockerRemove(ids, "rm")
		fmt.Println("Removed containers:", ids)

	case "removeimages":
		ids := getDockerIDs("images", "-q")
		runDockerRemove(ids, "rmi")

	case "listcontainers":
		runCommand("docker", "ps", "-a")

	case "listimages":
		runCommand("docker", "images")

	default:
		fmt.Println("Unknown command.")
	}
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
