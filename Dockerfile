FROM debian:bullseye-slim AS downloader

RUN apt-get update && apt-get install xz-utils -y

ADD https://github.com/yt-dlp/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz /usr/local/src/
RUN tar xf /usr/local/src/ffmpeg-master-latest-linux64-gpl.tar.xz -C /usr/local


FROM golang:1.20-alpine AS builder

WORKDIR /go/src/go.skymeyer.dev/media-manager/

COPY ./app/ ./app/
COPY ./pkg/ ./pkg/
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN CGO_ENABLED=0 go build -a -o mm app/main.go


FROM python:3.10.19-slim

COPY --from=downloader /usr/local/ffmpeg-master-latest-linux64-gpl/bin/ffmpeg /bin/ffmpeg
COPY --from=downloader /usr/local/ffmpeg-master-latest-linux64-gpl/bin/ffplay /bin/ffplay

RUN pip install --upgrade pip

RUN adduser runner
USER runner
WORKDIR /home/runner

RUN pip install --user yt-dlp
RUN pip install --user scdl

ENV PATH="/home/runner/.local/bin:${PATH}"

COPY --from=builder /go/src/go.skymeyer.dev/media-manager/mm /bin/mm

RUN mkdir /home/runner/download
WORKDIR /home/runner/download

ENTRYPOINT [ "mm" ]
