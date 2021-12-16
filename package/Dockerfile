FROM alpine
RUN apk -U add bash vim less curl wget
COPY bin/portexporter-linux /usr/bin/portexporter
ENTRYPOINT ["portexporter"]
