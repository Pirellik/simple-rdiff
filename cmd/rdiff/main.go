package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Pirellik/simple-rdiff/librsync"
)

const (
	helpCmd      string = "help"
	signatureCmd string = "signature"
	deltaCmd     string = "delta"
	patchCmd     string = "patch"

	helpMsg string = `Usage:
	rdiff help
	rdiff [options] signature old-file signature-file
	rdiff [options] delta signature-file new-file delta-file
	rdiff [options] patch basis-file delta-file new-file
Options:
	--block-size	size of the block in bytes
	`
)

type command interface {
	execute() error
}

type commandSignature struct {
	baseFilePath      string
	signatureFilePath string
	blockLength       uint32
}

func (c *commandSignature) execute() error {
	base, err := os.Open(c.baseFilePath)
	if err != nil {
		return err
	}
	defer base.Close()
	sig, err := librsync.NewSignature(base, c.blockLength)
	if err != nil {
		return err
	}
	sigFile, err := os.Create(c.signatureFilePath)
	if err != nil {
		return err
	}
	defer sigFile.Close()
	return sig.Write(sigFile)
}

type commandDelta struct {
	srcFilePath       string
	signatureFilePath string
	deltaFilePath     string
}

func (c *commandDelta) execute() error {
	src, err := os.Open(c.srcFilePath)
	if err != nil {
		return err
	}
	defer src.Close()
	sigFile, err := os.Open(c.signatureFilePath)
	if err != nil {
		return err
	}
	defer sigFile.Close()
	sig, err := librsync.ReadSignature(sigFile)
	if err != nil {
		return err
	}
	delta, err := librsync.NewDelta(src, sig)
	if err != nil {
		return err
	}
	deltaFile, err := os.Create(c.deltaFilePath)
	if err != nil {
		return err
	}
	defer deltaFile.Close()
	return delta.Write(deltaFile)
}

type commandPatch struct {
	baseFilePath  string
	deltaFilePath string
	outFilePath   string
}

func (c *commandPatch) execute() error {
	base, err := os.Open(c.baseFilePath)
	if err != nil {
		return err
	}
	defer base.Close()
	deltaFile, err := os.Open(c.deltaFilePath)
	if err != nil {
		return err
	}
	defer deltaFile.Close()
	delta, err := librsync.ReadDelta(deltaFile)
	if err != nil {
		return err
	}
	out, err := os.Create(c.outFilePath)
	if err != nil {
		return err
	}
	defer out.Close()
	return delta.Patch(base, out)
}

type commandHelp struct{}

func (c *commandHelp) execute() error {
	fmt.Println(helpMsg)
	return nil
}

func parseCmd() (command, error) {
	blockSize := flag.Int("block-size", 32, "size of the block in bytes")
	flag.Parse()
	values := flag.Args()
	if len(values) == 0 {
		return nil, errors.New("no command specified")
	}
	switch values[0] {
	case signatureCmd:
		if len(values) != 3 {
			return nil, errors.New("invalid signature command")
		}
		return &commandSignature{
			baseFilePath:      values[1],
			signatureFilePath: values[2],
			blockLength:       uint32(*blockSize),
		}, nil
	case deltaCmd:
		if len(values) != 4 {
			return nil, errors.New("invalid delta command")
		}
		return &commandDelta{
			signatureFilePath: values[1],
			srcFilePath:       values[2],
			deltaFilePath:     values[3],
		}, nil
	case patchCmd:
		if len(values) != 4 {
			return nil, errors.New("invalid patch command")
		}
		return &commandPatch{
			baseFilePath:  values[1],
			deltaFilePath: values[2],
			outFilePath:   values[3],
		}, nil
	case helpCmd:
		return &commandHelp{}, nil
	default:
		return nil, fmt.Errorf("invalid command: %s", values[0])
	}
}

func main() {
	cmd, err := parseCmd()
	if err != nil {
		fmt.Println(err)
		fmt.Println(helpMsg)
		os.Exit(1)
	}
	if err := cmd.execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
