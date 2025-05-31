# MANGO
CLI todo application for forgetful people who use the terminal

## Instalation
### Prerequisite
- go 1.24.3
### Using binaries
- Put the downloaded binary where $PATH points to
```bash
mv mango ~/.local/bin/mango # for example
```
### From sources
- Clone the repository
- Compile & install the program
```bash
go install
```
- Alternatively, you can compile, then put it anywhere where $PATH points to
```bash
go compile
mv mango ~/.local/bin/mango # for example
```
### Add it to the terminal
To have the todos pop up when starting the terminal, append this to the RC file of your shell
```bash
mango list
mango list --urgent     # Only urgent todos
mango list --urgent -n 5         # Limits to 5 todos
```
## Todo data
The todo data is stored in `~/.config/mango/todos.json`