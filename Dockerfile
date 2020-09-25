# Build restic
FROM golang:1.15-alpine as build
ADD . /app
WORKDIR /app
RUN go run build.go

# Get rclone
FROM alpine:latest as rclone
RUN apk add --no-cache curl unzip && \
	curl -s -o /tmp/rclone.zip https://downloads.rclone.org/v1.53.1/rclone-v1.53.1-linux-amd64.zip && \
	unzip -o /tmp/rclone.zip && \
	mv /rclone-*/rclone / && \
	rm -rf /rclone-* /tmp/rclone.zip

# Slap it all in an image
FROM alpine:latest
COPY --from=build /app/restic /usr/bin/restic
COPY --from=rclone /rclone /usr/bin/rclone
ADD aarnet/docker/ /
RUN apk add --no-cache bash && \
	chmod +x /usr/bin/rclone /usr/bin/restic /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]