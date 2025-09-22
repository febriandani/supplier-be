FROM golang:1.19

# Install sudo
RUN apt-get update && \
    apt-get install -y sudo

# Add a non-root user with sudo privileges
RUN useradd -ms /bin/bash user && \
    echo "user ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/app ./cmd/main.go

# Set the user
USER user

CMD ["/out/app"]

EXPOSE 8080
