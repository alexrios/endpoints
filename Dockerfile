FROM scratch
COPY endpoints /
ENTRYPOINT ["/endpoints"]
CMD ["--config=."]
