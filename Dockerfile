FROM golang:1.17.10-bullseye

WORKDIR /app
COPY . /app/

RUN make linux

EXPOSE 8080

CMD ["/app/dist/linux_amd64/mirage","serve"]
