# installer
Installer based on Kardianos Service

# Go Service with Interval Execution

This is a Go-based service that periodically executes a batch file at a specified interval. It is based on the original work by Daniel Theophanes (kardianos/service) and has been modified to support interval-based execution of a batch script.

## Features

- Runs a specified batch file at regular intervals.
- Configurable through a JSON file.
- Logs standard output and error to specified files.

## Configuration

Create a JSON configuration file in the same directory as the service executable. The configuration file should have the same name as the executable but with a `.json` extension (e.g., `myservice.json`). Here is an example configuration:
```
{
    "Name": "MyService",
    "DisplayName": "My Custom Service",
    "Description": "This service runs a batch file at regular intervals.",
    "Dir": "C:\\programming\\Go\\installer",
    "Exec": "C:\\programming\\Go\\installer\\script.bat",
    "Args": [],
    "Env": [],
    "Stderr": "service_error.log",
    "Stdout": "service_output.log",
    "Interval": 1  // Interval in minutes
}
```
## Usage

### Prerequisites

- Go 1.20 or later

### Download the Source Code

1. Clone the repository:

git clone https://github.com/AASoftware/installer.git
cd installer

### Compile the Executable

2. Build the executable:

go build -o myservice.exe

3. (Optional) Obfuscate the build:

go install mvdan.cc/garble@latest
garble build -o myservice.exe

### Install the Service

4. Install the service:

myservice.exe -service install

5. Start the service:

myservice.exe -service start

### Stop and Uninstall the Service

- To stop the service:

myservice.exe -service stop

- To uninstall the service:

myservice.exe -service uninstall

## Acknowledgements

This project is based on the work of Daniel Theophanes (kardianos/service) and has been modified to include interval-based execution of a batch script.

## License

This project is licensed under the terms of the zlib License (LICENSE).
