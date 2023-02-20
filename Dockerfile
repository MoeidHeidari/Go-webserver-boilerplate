# Use an official Go image as the base image
FROM golang:1.18.1

# Set the working directory to /app
WORKDIR /.

# Copy the public.key file to the container
COPY public.key /app/public.key

# Set environment variables
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# Copy the current directory contents into the container at /app
COPY . .
RUN go get -d -v 
# Expose ports 3000 and 12121
EXPOSE 3000
EXPOSE 12121

# Set the command to run the Go app
CMD ["go", "run", ".", "app:serve"]