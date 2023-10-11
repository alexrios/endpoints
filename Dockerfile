FROM alpine:3.18
COPY endpoints /
ENTRYPOINT ["/endpoints"]
CMD ["--config=."]
