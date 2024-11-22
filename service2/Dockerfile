# Use an official Go runtime as a parent image
FROM golang:1.23.3

ENV export GOARM=7
# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any needed packages specified in requirements.txt
RUN go mod download

# Run app.py when the container launches
CMD ["go", "run", "service2.go"]