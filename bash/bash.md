## Bash
### syntax checker - shellcheck
https://www.shellcheck.net/
### tests(hopefully I can use a real language) - bats
https://github.com/sstephenson/bats
### CLI - simplified man page - tldr
https://github.com/tldr-pages/tldr
### CLI - faster navigation of directories - autojump
https://www.linode.com/docs/tools-reference/tools/faster-file-navigation-with-autojump/
### CLI - Html preview in the terminal - w3m
```bash
sudo yum install w3m -y
w3m -dump /some/path/file.html
```
Quite usefull to visualize html based report, like checkstyle, etc, on a remote host
### pyenv
See https://github.com/pyenv/pyenv  
and its installer: https://github.com/pyenv/pyenv-installer
### NVM
See https://github.com/nvm-sh/nvm  
Installation instructions here: https://github.com/nvm-sh/nvm#installation-and-update
### RVM - https://rvm.io/
```bash
gpg2 --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB
curl -sSL https://get.rvm.io | bash -s stable
```
### Run previous command with sudo
```bash
alias please="sudo $(fc -ln -1)"
```
### Get the source directory of a bash script from within the script itself
```bash
my_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
```
It will work as long as the last component of the path used to find the script 
is not a symlink (directory links are OK).
### CLI - Simplistic interactive filtering tool
https://github.com/peco/peco  
__Example:__ Search file, select it and cat it
```bash
find -name '*.java' | peco | xargs cat
```




