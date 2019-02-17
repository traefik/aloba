FROM golang:1-alpine as builder

RUN apk --update upgrade \
&& apk --no-cache --no-progress add git make \
&& rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/containous/aloba
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN make dependencies
RUN make build

FROM alpine:3.6
RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates \
    && rm -rf /var/cache/apk/*

LABEL "com.github.actions.name"="Aloba"
LABEL "com.github.actions.description"="Add labels and milestone on pull requests and issues"
LABEL "com.github.actions.icon"="cpu"
LABEL "com.github.actions.color"="purple"

LABEL "repository"="http://github.com/containous/aloba"
LABEL "homepage"="http://github.com/containous/aloba"
LABEL "maintainer"="ldez <ldez@users.noreply.github.com>"

COPY --from=builder /go/src/github.com/containous/aloba/aloba /usr/bin/aloba

ENTRYPOINT ["/usr/bin/aloba"]
