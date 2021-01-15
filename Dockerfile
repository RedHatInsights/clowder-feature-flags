FROM golang as compiler 
RUN CGO_ENABLED=0 go get -a -ldflags '-s' \ 
github.com/redhatinsights/clowder-feature-flags

FROM scratch 
COPY --from=compiler /go/bin/clowder-feature-flags . 
CMD ["./clowder-feature-flags"]