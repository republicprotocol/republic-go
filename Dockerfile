FROM golang:1.8

RUN mkdir -p $GOPATH/src/github.com/lambdaprotocol/lambda
WORKDIR $GOPATH/src/github.com/lambdaprotocol/lambda
COPY . .

RUN go install

EXPOSE 5000

CMD ["lambda"]