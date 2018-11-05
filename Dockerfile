FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY fileserver /server
EXPOSE 1212
ENTRYPOINT ["/server"]
