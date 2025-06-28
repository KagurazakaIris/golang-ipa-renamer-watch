# golang-ipa-renamer-watch

A simple Go tool to watch a directory for new IPA files and automatically rename them based on your rules.

## Features
- Watches a specified directory for new or changed IPA files
- Automatically renames IPA files according to custom rules
- Easy to configure and use
- Docker support for easy deployment

## Usage

### Build and Run Locally

1. Clone the repository:
   ```sh
   git clone https://github.com/<your-username>/golang-ipa-renamer-watch.git
   cd golang-ipa-renamer-watch
   ```
2. Build the project:
   ```sh
   go build -o ipa_renamer main.go
   ```
3. Run the tool:
   ```sh
   ./ipa_renamer
   ```

### Run with Docker

1. Build the Docker image:
   ```sh
   docker build -t ipa-renamer .
   ```
2. Run the container:
   ```sh
   docker run -v /path/to/watch:/watch ipa-renamer
   ```

### Run with GitHub Container Registry (GHCR)

1. Pull the image:
   ```sh
   docker pull ghcr.io/<your-username>/golang-ipa-renamer-watch:latest
   ```
2. Run the container:
   ```sh
   docker run -v /path/to/watch:/watch ghcr.io/<your-username>/golang-ipa-renamer-watch:latest
   ```

## Configuration

You can configure the watch directory and renaming rules via environment variables or command-line arguments (see `main.go` for details).

## License

This project is licensed under the MIT License.

