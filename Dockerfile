FROM scratch

WORKDIR /
ADD tag-exporter /
ENTRYPOINT ["/tag-exporter"]