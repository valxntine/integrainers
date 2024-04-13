# syntax=docker/dockerfile:1.0.0-experimental

# BUILD STAGE - Stage for building the app
#
FROM golang:1.21-bullseye
#
#RUN mkdir -p -m 0600 /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts
#RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

WORKDIR /booky

COPY go.mod go.sum Makefile ./

RUN make deps

COPY . /booky

RUN go build -o /book-service

EXPOSE 8088

CMD ["/book-service"]
