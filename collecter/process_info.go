package collecter

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

func GetCPUAndMemChrome() (map[string]float64, error) {
	outputMap := map[string]float64{"cpu": 0.0, "mem": 0.0, "num": 0.0}
	cmd := exec.Command("/bin/bash", "-c", `ps -e -o 'comm,pcpu,rsz'`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return nil, err
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return nil, err
	}

	for _, x := range strings.Split(string(bytes), "\n") {
		if len(strings.Fields(x)) != 0 && strings.Fields(x)[0] == "chrome" {
			cpuTemp, _ := strconv.ParseFloat(strings.Fields(x)[1], 64)
			outputMap["cpu"] += cpuTemp
			memTemp, _ := strconv.ParseFloat(strings.Fields(x)[2], 64)
			outputMap["mem"] += memTemp
			outputMap["num"] += 1.0
		}
	}

	return outputMap, nil
}

func IsChromeClosed() bool {
	cmd := exec.Command("/bin/bash", "-c", `ps -e -o 'comm,pcpu,rsz'`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
	}

	for _, x := range strings.Split(string(bytes), "\n") {
		if len(strings.Fields(x)) != 0 && strings.Fields(x)[0] == "chrome" {
			return false
		}
	}
	return true
}
