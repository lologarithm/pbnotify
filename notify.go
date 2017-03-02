package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] != "songstart" {
			os.Exit(0)
		}
	}

	songinfo := map[string]string{}

	bio := bufio.NewReader(os.Stdin)
	var err error
	for ; err == nil; _, err = bio.Peek(1) {
		line, err := bio.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Trim(strings.TrimSpace(line), "\n") // clean up string
		kv := strings.Split(line, "=")
		songinfo[kv[0]] = kv[1]
	}
	cmd := exec.Command("notify-send", fmt.Sprintf("%s by %s", songinfo["title"], songinfo["artist"]), "Now Playing")
	if songinfo["coverArt"] != "" {
		if art := getArt(songinfo["coverArt"]); art != nil {
			cmd.Args = append(cmd.Args, "-i")
			cmd.Args = append(cmd.Args, art.Name())
			defer os.Remove(art.Name())
		}
	}

	cmd.Run()
}

func getArt(uri string) *os.File {
	resp, err := http.Get(uri)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	f, err := ioutil.TempFile("", "pbnot")
	if err != nil {
		fmt.Printf("Failed to create temp file: %s\n", err)
		return nil
	}
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		return nil
	}

	return f
}
