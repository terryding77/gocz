# gocz

gocz is a go version of commitizen

## install and usage
1. use `go get -u -v github.com/terryding77/gocz/cli/git-cz` download and install binary file
2. use `ls -al $GOPATH/bin/git-cz` ensure binary file created
3. keep `$GOPATH/bin` in your **PATH** environment variable, if not, use `export PATH=${PATH}:$GOPATH/bin`
4. you can copy the basic .gocz.toml file in $GOPATH/src/github.com/terryding77/gocz/sample/ to your git repo root, or home root.
5. use operation `git cz` to replace `git commit`

