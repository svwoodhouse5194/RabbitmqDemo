# Use an official Python runtime as a parent image
FROM golang:alpine3.20

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any needed packages specified in requirements.txt
RUN go mod download

# Run app.py when the container launches
CMD ["go", "run", "service1.go"]