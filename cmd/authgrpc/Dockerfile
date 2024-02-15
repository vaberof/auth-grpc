FROM golang:1.22

RUN mkdir -p /opt/app

COPY cmd/authgrpc/build /opt/app

WORKDIR /opt/app

CMD [ "/opt/app/main", "-config.files", "container.yaml", "-env.vars.file", "sample.env" ]