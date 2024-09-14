################################
# STEP 1 build executable binary
################################
FROM --platform=${BUILDPLATFORM} golang:1.23.1 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Set the current working directory inside the container
WORKDIR /opt/flawlessbyte/nhs-bank-notifier

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the binary.
# flags:
# -s strips the binary by omitting the symbol table and debug information
# -w further strips the binary by also omitting the DWARF symbol table
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-s -w" \
    -o /usr/local/bin/nhs-bank-notifier \
    ./cmd/nhs-bank-notifier/main.go

########################################
# STEP 2 build a small image from alpine
########################################
FROM --platform=${TARGETPLATFORM} alpine
RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/local/bin/nhs-bank-notifier /usr/local/bin/nhs-bank-notifier

EXPOSE 8080
# Command to run the executable
ENTRYPOINT ["/usr/local/bin/nhs-bank-notifier"]