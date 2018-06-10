# build stage
FROM golang:1.9-alpine AS build-env

# ARG SSH_PRIVATE_KEY
# ARG SSH_PUB_KEY
RUN apk update && apk add git
#&& apk add mercurial && apk add openssh-client

# MICRO DEPENDENCY
RUN go get github.com/micro/micro/cmd
# RUN go get github.com/micro/micro/cli
# # Authorize SSH Host
# RUN mkdir -p /root/.ssh && \
#     chmod 0700 /root/.ssh && \
#     touch /root/.ssh/known_hosts && \
#     ssh-keyscan github.com > /root/.ssh/known_hosts && \
#     ssh-keyscan bitbucket.org > /root/.ssh/known_hosts

# # Add the keys and set permissions
# RUN echo $SSH_PRIVATE_KEY > /root/.ssh/id_rsa && \
#     echo $SSH_PUB_KEY > /root/.ssh/id_rsa.pub && \
#     chmod 600 /root/.ssh/id_rsa && \
#     chmod 600 /root/.ssh/id_rsa.pub

# RUN mkdir -p /go/src \
#   && mkdir -p /go/bin \
#   && mkdir -p /go/pkg \
#   && mkdir -p /go/src/github.org/micro/micro/

# # Get Travelplatform dependencies
# RUN git config --global url."git@bitbucket.org:".insteadOf "https://api.bitbucket.org/"
# WORKDIR go/src/bitbucket.org/appgoplaces
# RUN git clone git@bitbucket.org:appgoplaces/travelplatform-system/lib/plugins/auth

WORKDIR /go/src/github.org/micro/micro
ADD . .
RUN go build -v -o micro ./main.go ./plugin.go

# final stage
FROM alpine:3.2
RUN apk add --update ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.org/micro/micro/micro micro
ENTRYPOINT [ "/app/micro" ]
