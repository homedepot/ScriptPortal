from golang:alpine
RUN apk update
RUN apk add git gcc musl-dev bash python3
RUN go get github.com/homedepot/ScriptPortal
WORKDIR /go/src/github.com/homedepot/ScriptPortal
RUN ./build.sh
RUN ./install.sh
CMD ScriptPortal
