# syntax=docker/dockerfile:1.0.0-experimental

# BUILD STAGE - Stage for building the app
#
FROM golang:1.21-bullseye as base_build_app
#
#RUN mkdir -p -m 0600 /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts
#RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

WORKDIR /booky/

COPY go.mod go.sum Makefile ./

RUN make deps

FROM base_build_app as build_app

COPY . /booky

RUN make build

# TEST STAGE - Stage for testing, includes dev dependencies
#
FROM base_build_app as test

COPY --from=build_app /booky/ .
