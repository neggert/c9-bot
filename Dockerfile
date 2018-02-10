FROM alpine:3.7
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD main /
CMD ["/main"]