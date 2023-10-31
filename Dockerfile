FROM golang:1.21 as builder

WORKDIR /src

COPY . .

ENV USER=dump1090
ENV UID=1001

ARG BUILDID
ARG RUNENV

RUN adduser --disabled-password --gecos "" --home "/tmp" --shell "/sbin/nologin" --no-create-home --uid "${UID}" "${USER}" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-w -s" -o dump1090-exporter . 

FROM ubuntu:20.04

WORKDIR /app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=1001:1001 /src/dump1090-exporter /app/dump1090-exporter

USER dump1090:dump1090

EXPOSE 9467

CMD [ "/app/dump1090-exporter" ]
