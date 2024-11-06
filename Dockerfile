FROM gcr.io/distroless/static AS distroless
FROM scratch
COPY --from=distroless /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY server /server
EXPOSE 5051
ENTRYPOINT ["/server"]
