package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	defaultDataDir, err := getDefaultDataDir()
	if err != nil {
		return err
	}

	app := &cli.App{
		Name:   "scratch",
		Usage:  "Quick terminal notes",
		Action: handleEditTmpFile,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "editor",
				Aliases: []string{"e"},
				Value:   "vim",
				Usage:   "A terminal command to open the desired text editor",
			},
			&cli.StringFlag{
				Name:  "data-dir",
				Value: defaultDataDir,
				Usage: "The directory to use for story scratch pads",
			},
			&cli.StringFlag{
				Name:    "ext",
				Aliases: []string{"x"},
				Value:   "md",
				Usage:   "Specify the file extension",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "pad",
				Usage:  "Edit the shared persistant scratch pad",
				Action: handleEditScratchPad,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "fresh",
						Usage: "Open a fresh scratch pad (this irreversibly deletes the old data)",
						Value: false,
					},
					&cli.StringFlag{
						Name:    "ext",
						Aliases: []string{"x"},
						Value:   "md",
						Usage:   "Specify the file extension",
					},
				},
			},
			{
				Name:   "add",
				Usage:  "Create a new scratch pad",
				Action: handleAddPad,
			},
			{
				Name:   "ls",
				Usage:  "List all pads",
				Action: handleListPads,
			},
			{
				Name:   "edit",
				Usage:  "Edit a pad",
				Action: handleEditPad,
			},
			{
				Name:   "rm",
				Usage:  "Remove a pad",
				Action: handleRemovePad,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "recursive",
						Aliases: []string{"r"},
						Usage:   "Recursively delete a directory",
					},
				},
			},
		},
		Before: func(c *cli.Context) error {
			dataDir := c.String("data-dir")
			if dataDir == "" {
				return errors.New("data-dir cannot be empty")
			}

			// Before performing any operations, ensure that all required directories exist
			if err := createDirIfNotExists(filepath.Join(dataDir, "pads")); err != nil {
				return err
			}

			if err := createDirIfNotExists(filepath.Join(dataDir, "defaults")); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

func getDefaultDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".scratch"), nil
}

func createDirIfNotExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func buildFileName(filename, ext string) string {
	if ext == "" {
		return filename
	}

	return fmt.Sprintf("%s.%s", filename, strings.TrimLeft(ext, "."))
}

func handleEditTmpFile(c *cli.Context) error {
	if c.NArg() > 0 {
		return fmt.Errorf("unkown command: %s", c.Args().First())
	}

	f, err := os.CreateTemp("", buildFileName("scratch-*", c.String("ext")))
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	return editFile(c.String("editor"), f.Name())
}

func handleEditScratchPad(c *cli.Context) error {
	filename := filepath.Join(
		c.String("data-dir"),
		"defaults",
		buildFileName("scratch", c.String("ext")),
	)

	if c.Bool("fresh") {
		err := os.WriteFile(filename, []byte{}, 0755)
		if err != nil {
			return fmt.Errorf("unable to clear old scratch pad: %v", err)
		}
	}

	return editFile(c.String("editor"), filename)
}

func handleAddPad(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("pad name(s) required")
	}

	for _, pad := range c.Args().Slice() {
		if err := addPad(c.String("data-dir"), pad); err != nil {
			return err
		}
	}

	return nil
}

func handleListPads(c *cli.Context) error {
	entries, err := os.ReadDir(filepath.Join(c.String("data-dir"), "pads", c.Args().First()))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}

		fmt.Println(name)
	}

	return nil
}

func handleEditPad(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("pad name required")
	}

	padName := c.Args().First()

	if path := filepath.Dir(padName); path != "." {
		err := createDirIfNotExists(filepath.Join(c.String("data-dir"), "pads", path))
		if err != nil {
			return err
		}
	}

	filename := filepath.Join(c.String("data-dir"), "pads", padName)
	return editFile(c.String("editor"), filename)
}

func handleRemovePad(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("pad name required")
	}

	path := filepath.Join(c.String("data-dir"), "pads", c.Args().First())

	f, err := os.Stat(path)
	if err != nil {
		return err
	}

	if f.IsDir() && !c.Bool("recursive") {
		return errors.New("must recursively delete directories")
	}

	return os.RemoveAll(path)
}

func addPad(dataDir, padName string) error {
	if path := filepath.Dir(padName); path != "." {
		err := createDirIfNotExists(filepath.Join(dataDir, "pads", path))
		if err != nil {
			return err
		}
	}

	f, err := os.Create(filepath.Join(dataDir, "pads", padName))
	if err != nil {
		return err
	}

	return f.Close()
}

func editFile(editor, filename string) error {
	editorExec, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	c := exec.Command(editorExec, filename)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
