FROM golang:1-alpine as builder

RUN apk --update upgrade \
&& apk --no-cache --no-progress add git make \
&& rm -rf /var/cache/apk/*

WORKDIR /go/aloba

ENV GO111MODULE on

# Download go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.10
RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates \
    && rm -rf /var/cache/apk/*

LABEL "com.github.actions.name"="Aloba"
LABEL "com.github.actions.description"="Add labels and milestone on pull requests and issues"
LABEL "com.github.actions.icon"="cpu"
LABEL "com.github.actions.color"="purple"

LABEL "repository"="http://github.com/traefik/aloba"
LABEL "homepage"="http://github.com/traefik/aloba"
LABEL "maintainer"="ldez <ldez@users.noreply.github.com>"

COPY --from=builder /go/aloba/aloba /usr/bin/aloba

ENTRYPOINT ["/usr/bin/aloba"]
