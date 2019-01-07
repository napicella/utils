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
### RVM - https://rvm.io/
```bash
gpg2 --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB
curl -sSL https://get.rvm.io | bash -s stable
```
### Run previous command with sudo
```bash
alias please="sudo $(fc -ln -1)"
```


