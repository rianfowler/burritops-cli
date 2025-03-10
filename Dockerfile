FROM buildpacksio/pack:latest

# Set the working directory where your code will be mounted.
WORKDIR /workspace

# Default command: build an app called "my-app" using a specified builder.
# Adjust the builder image and app name as necessary.
CMD ["pack", "build", "my-app", "--builder", "paketobuildpacks/builder:base", "--path", "."]

