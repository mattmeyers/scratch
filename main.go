package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	defaultDataDir, err := getDefaultDataDir()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:   "scratch",
		Usage:  "quick notes",
		Action: handleEditScratchPad,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "editor",
				Value:   "vim",
				Aliases: []string{"e"},
			},
			&cli.StringFlag{
				Name:  "data-dir",
				Value: defaultDataDir,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getDefaultDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "scratch"), nil
}

func createDirIfNotExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func handleEditScratchPad(c *cli.Context) error {
	if err := createDirIfNotExists(c.String("data-dir")); err != nil {
		return err
	}

	return editFile(c.String("editor"), filepath.Join(c.String("data-dir"), "scratch.md"))
}

func editFile(editor, filename string) error {
	editorExec, err := exec.LookPath(editor)
	if err != nil {
		log.Fatal(err)
	}

	c := exec.Command(editorExec, filename)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
