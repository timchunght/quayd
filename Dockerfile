FROM golang:1.8.3

WORKDIR /go/src/github.com/timchunght/quayd
COPY . .

RUN cd cmd/quayd && go install
CMD ["quayd", "-github-token=$GITHUB_TOKEN", "-port=$PORT", "-registry-auth=$REGISTRY_AUTH"]