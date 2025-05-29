# fyora
> The Faerie Queen, Fyora, is the ruler of Faerieland. She is a kind faerie who basically tries to keep everything under control, not just in Faerieland. (Neopets Wiki)

Fyora is a declarative replacement for GNU Stow. It is a symlink farm manager that uses a yaml file to declare which files/directories should be symlinked where, rather than a cli interface.

## Usage
1. Create a config file (see below)
2. Run `fyora`
    1. By default, it will look for a file called `fyora.yaml` in your `.config` directory. You can specify a different path using the `--config` flag (`-c` shorthand).

## Config
(the actually important part)

The fyora config file is a yaml file that looks like this:
```yaml
links:
    - type: outside
      source: /dir1
      target: ~/dir2
    - type: inside
      source: ~/dir3
      target: ~/dir/dir4
    - type: file
      source: /dir5/file.txt
      target: ~/dir2/dir/file.txt
ignore:
    - .DS_Store
    - .git
```
Outside links create a symlink to the folder itself i.e. ~/dir2/dir1 be symlinked to /dir1.

Inside links symlink everything inside the first folder inside of the second folder i.e. ~/dir3/file.txt would be symlinked to ~/dir/dir4/file.txt.

File links symlink the file itself i.e. /dir5/file.txt would be symlinked to ~/dir2/dir/file.txt.

Everything under ignore (files and folders) will NOT be symlinked.

## Installation
1. Clone the repo
2. Run `go mod tidy` to install dependencies
3. Run `go build` to build the binary
4. Run `go install` to install the binary to your $GOPATH/bin directory
5. Make a `fyora.yaml` file in your home directory (or wherever you want) and add your links to it
6. Run `fyora` to create the symlinks (in the same directory as `fyora.yaml`, this is being worked on)

## Where does the name come from?
> "The Faerie Queen, Fyora, is the ruler of Faerieland. She is a kind faerie who basically tries to keep everything under control, not just in Faerieland." (Neopets Wiki)

My friend is a neopets fan and really wanted me to name it after Fyora. I thought it was a cute name and it stuck.
