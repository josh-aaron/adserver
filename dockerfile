FROM golang:1.25 AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# RUN if [ "$TARGET_OS" = "windows" ]; then \
#     echo "Running windows build commands"; \ 
#     cd cmd/api; \
# 	go build -o ../../bin/adserver.exe;\ 
#     else \
#     echo "Running mac build commands"; \ 
#     go build -o adserver cmd/api/*.go; \
#     fi

EXPOSE 8000

CMD ["/build/adserver"]

