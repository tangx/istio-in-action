FROM alpine
ADD out/prod /go/bin/prod
EXPOSE 80
ENTRYPOINT ["/go/bin/prod"]

