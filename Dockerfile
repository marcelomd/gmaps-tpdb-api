FROM golang

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o api ./cmd/api
EXPOSE 8080
CMD [ "/app/api" ]
