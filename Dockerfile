FROM golang:1.19.1-buster as go-target
ADD . /grate
WORKDIR /grate/cmd/grater
RUN go build

FROM golang:1.19.1-buster
COPY --from=go-target /grate/cmd/grater/grater /
COPY --from=go-target /grate/testdata/*.xlsx /testsuite/

ENTRYPOINT []
CMD ["/grater", "@@"]
