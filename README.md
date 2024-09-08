# antares

antares is user-friendly online interface, Antares is a lightweight, cross-platform web drive solution that makes managing and accessing files simple. It enables you to rapidly set up a file server on Windows, Linux, or macOS so that you may access your files online!

## Features

- Easy file management through web interface
- Cross-platform support (Windows, Linux, macOS)
- File operations: upload, download, delete, rename, move
- Directory creation and management
- Optional request logging

## Usage

### Basic Usage
To start Antares with default settings:
```
antares
```

This will start the server on the default port (9700) and serve the current directory.

### Command-line Options

Antares supports several command-line options:

* -gen: Generate a default configuration file
* -port <int>: Specify the port to run the server on
* -root <string>: Set the root directory to serve
* -log: Enable request logging
* -open: Open the default browser after starting the server

## License

This project is licensed under The MIT.
