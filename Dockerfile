FROM golang:1.17.10-bullseye

WORKDIR /app
COPY . /app/

RUN make linux

EXPOSE 6456

CMD ["/app/dist/linux_amd64/mirage","serve"]
