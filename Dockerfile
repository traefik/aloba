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
COPY --from=builder /go/src/github.com/containous/aloba/aloba .
CMD ["./aloba", "-h"]
