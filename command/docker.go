package command

import (
	"log"
	"os/exec"
)

func DockerUp() {
	rcmd := `docker compose up -d`
	cmd := exec.Command("bash", "-c", rcmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(string(out))
}

func DockerStop() {
	rcmd := `docker compose stop`
	cmd := exec.Command("bash", "-c", rcmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(string(out))
}
