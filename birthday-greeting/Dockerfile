FROM public.ecr.aws/bitnami/golang:1.16 as build-image

# Don't change this folder!
WORKDIR /go/src

# complained by SAM: properly copy the go environment into docker
# TODO: copy pkg to save compile time
COPY go.mod go.sum main.go ./
COPY dao ./dao
COPY types ./types
COPY utils ./utils

RUN export GO111MODULE="on"
RUN go build ./...
RUN go mod tidy
RUN go mod download

RUN go build -o ../bin

FROM public.ecr.aws/lambda/go:1

COPY --from=build-image /go/bin/ /var/task/

# Command can be overwritten by providing a different command in the template directly.
CMD ["birthday-greeting"]
