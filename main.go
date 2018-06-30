package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gobuffalo/packr"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func exec(minS string) (err error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	min, err := strconv.Atoi(minS)
	if err != nil {
		return fmt.Errorf("failed parse min, it should be int: %s: %v", minS, err)
	}
	sec := min * 60
	s.Start()
	for i := 0; i < sec; i++ {
		s.Suffix = fmt.Sprintf(" %s", sec2str(sec-i))
		time.Sleep(1 * time.Second)
	}
	s.Stop()

	box := packr.NewBox("./sound")

	// https://maoudamashii.jokersounds.com/archives/se_maoudamashii_onepoint23.html
	rc := ioutil.NopCloser(bytes.NewReader(box.Bytes("se_maoudamashii_onepoint23.mp3")))
	d, err := mp3.NewDecoder(rc)
	if err != nil {
		return err
	}
	defer close(d)

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer close(p)

	_, err = io.Copy(p, d)
	if err != nil {
		return err
	}
	fmt.Printf("Ended %s\n", sec2str(sec))
	return nil
}

func sec2str(sec int) string {
	return fmt.Sprintf("%02d:%02d", sec/60, sec%60)
}

func close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	mes := `ktimer is a simple timer on cli.

Usage:
	ktimer minutes
	ex. ktimer 30 -> start 30 minutes timer.

Args:
	minutes	time to set. acceptable only int.

LICENSE:
	This program provided under MIT.

CREDIT:
	Program: Yutaka Nishimura github.com/ynishi
	Sound: 魔王魂 https://maoudamashii.jokersounds.com/
`
	fmt.Print(mes)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}
	t := os.Args[1]
	if err := exec(t); err != nil {
		log.Fatalf("failed to exec: %v", err)
	}
}
