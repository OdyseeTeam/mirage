FROM golang:1.19-bullseye

WORKDIR /app
COPY . /app/
RUN apt-get update && apt-get install -y libvips libvips-dev
RUN make linux

EXPOSE 6456

CMD ["/app/dist/linux_amd64/mirage","serve"]
