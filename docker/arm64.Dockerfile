#################################################################
# arm64 전용 빌드 단계
FROM golang:1.24-alpine AS build-base
RUN apk add --no-cache zip make git tar
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ENV CGO_ENABLED=0
RUN rm -rf tmp binaries
RUN mkdir tmp binaries
RUN cp mediamtx.yml LICENSE tmp/
RUN go generate ./...

# arm64 바이너리 빌드
FROM build-base AS build-linux-arm64
ENV GOOS=linux GOARCH=arm64
RUN go build -o "tmp/mediamtx"
RUN tar -C tmp -czf "binaries/mediamtx_$(cat internal/core/VERSION)_linux_arm64.tar.gz" --owner=0 --group=0 "mediamtx" mediamtx.yml LICENSE

# 바이너리 추출 단계
FROM build-base AS extract-binaries
COPY --from=build-linux-arm64 /s/binaries/mediamtx_*_linux_arm64.tar.gz /tmp/
RUN cd /tmp && tar -xzf mediamtx_*_linux_arm64.tar.gz

#################################################################
# 최종 실행 이미지 (arm64만)
FROM --platform=linux/arm64 alpine:3.20

RUN apk add --no-cache ffmpeg python3 py3-pip py3-opencv py3-numpy

COPY --from=extract-binaries /tmp/mediamtx /mediamtx
COPY --from=extract-binaries /tmp/mediamtx.yml /mediamtx.yml
COPY --from=extract-binaries /tmp/LICENSE /LICENSE

ENTRYPOINT [ "/mediamtx" ] 