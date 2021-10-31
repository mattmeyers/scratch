# scratch

`scratch` is a simple terminal notes tool. It provides ephemeral scratch pads, quick access persistant pads, and structured directories.

## Usage

```
NAME:
   scratch - Quick terminal notes

USAGE:
   scratch [global options] command [command options] [arguments...]

COMMANDS:
   pad      Edit the shared persistant scratch pad
   add      Create a new scratch pad
   ls       List all pads
   edit     Edit a pad
   rm       Remove a pad
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --editor value, -e value  A terminal command to open the desired text editor (default: "vim")
   --data-dir value          The directory to use for story scratch pads (default: "/home/matt/scratch")
   --ext value, -x value     Specify the file extension (default: "md")
   --help, -h                show help (default: false)
```