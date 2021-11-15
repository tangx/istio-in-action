FROM alpine
ADD out/review /go/bin/review
EXPOSE 80
ENTRYPOINT ["/go/bin/review"]

