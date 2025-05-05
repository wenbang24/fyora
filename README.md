# fyora
> The Faerie Queen, Fyora, is the ruler of Faerieland. She is a kind faerie who basically tries to keep everything under control, not just in Faerieland. (Neopets Wiki)

Fyora is a declarative replacement for GNU Stow. It is a symlink farm manager that uses a yaml file to declare which files/directories should be symlinked where, rather than a cli interface.

## Usage
Literally just `fyora`.

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

## Where does the name come from?
My friend plays neopets and really wanted me to name it after Fyora, the faerie queen. I thought it was a cute name and it stuck.