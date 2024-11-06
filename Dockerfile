FROM debian
COPY server /server
ENTRYPOINT ["/server"]
