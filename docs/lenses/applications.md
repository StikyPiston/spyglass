# Configuring the Applications lens

Configuration takes place in the `~/.config/spyglass/applications` directory.

To create an application entry, create a `.yaml` file named after that entry.

```yaml
name: "LibreOffice"                      # The name of the application
description: "Create and edit documents" # The description
icon: ""                                # Icon, typically from a nerd font, or an emoji
command: "libreoffice & sleep 5"         # Command to run. Note that the `& sleep 5` is needed because sometimes the terminal will close before the program can be detached, and it won't launch
context:                                 # Context menu items
  - name: "Open Calc"                    # Name of the context action
    command: "libreoffice --calc"        # Command for the context action to run
  - name: "Open Impress"
    command: "libreoffice --impress"
  - name: "Open Writer"
    command: "libreoffice --writer"
```
