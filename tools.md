### tldr - CLI - simplified man page
https://github.com/tldr-pages/tldr

### autojump - CLI - faster navigation of directories - 
https://www.linode.com/docs/tools-reference/tools/faster-file-navigation-with-autojump/

### w3m - CLI - Html preview in the terminal
```bash
sudo yum install w3m -y
w3m -dump /some/path/file.html
```
Quite usefull to visualize html based report, like checkstyle, etc, on a remote host

### img2pdf - CLI - Images to pdf
https://gitlab.mister-muffin.de/josch/img2pdf
```
img2pdf *.jpf --output docs.pdf
```

### awslogs - CLI - Get AWS lambda logs in your terminal
https://github.com/jorgebastida/awslogs
```
awslogs get /aws/lambda/SomeFunction ALL --watch
```

### peco - CLI - Simplistic interactive filtering tool
https://github.com/peco/peco  
__Example:__ Search file, select it and cat it
```bash
find -name '*.java' | peco | xargs cat
```

### fzf - CLI - Interactive finder
https://github.com/junegunn/fzf#usage  
Similar to peco but more customazible. Can work as quick way to build interactive menu/selection in bash
```
printf "item-1\nitem2\nitem3\n" | fzf --height 5
```

### thefuck - CLI - Corrects errors in previous console commands
See https://github.com/nvbn/thefuck  
```bash
> ls /root/
ls: cannot open directory '/root/': Permission denied

> fuck
sudo ls --color=auto /root/ [enter/↑/↓/ctrl+c]
[sudo] password for user: 
```

### pyenv - Python version manager
See https://github.com/pyenv/pyenv  
and its installer: https://github.com/pyenv/pyenv-installer

### nvm - Node version manager
See https://github.com/nvm-sh/nvm  
Installation instructions here: https://github.com/nvm-sh/nvm#installation-and-update

### rvm - Ruby version manager
See https://rvm.io/
```bash
gpg2 --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB
curl -sSL https://get.rvm.io | bash -s stable
```
### sdkman - Java (and other) version manager
https://sdkman.io/

### peek - Simple animated GIF screen recorder 
https://github.com/phw/peek

### shellcheck - syntax checker
https://www.shellcheck.net/

### bats - Writing tests for bash scripts
https://github.com/sstephenson/bats

### Run previous command with sudo
```bash
alias please="sudo \$(fc -ln -1)"
```
### Get the source directory of a bash script from within the script itself
```bash
my_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
```
It will work as long as the last component of the path used to find the script 
is not a symlink (directory links are OK).

### Simple http server based on nc
```bash
#!/bin/bash
while true; do printf 'HTTP/1.1 200 OK\n\n%s' "$(cat index.html)" | nc -l 5555; done
```





