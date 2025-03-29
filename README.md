# DockerDB

DockerDB is a command-line utility designed to simplify the setup and management of Docker containers for various databases. This tool provides an easy way to create, start, and manage containers for popular databases such as MySQL, MariaDB, PostgreSQL, MongoDB, and Redis.

## Features

- **Easy Setup**: Quickly set up Docker containers for different databases with simple commands.
- **Multiple Database Support**: Supports MySQL, MariaDB, PostgreSQL, MongoDB, and Redis.
- **Docker API Integration**: Interact with the Docker API to manage containers seamlessly.

## Installation

To install DockerDB, clone the repository and build the project:

```bash
git clone <repository-url>
cd dockerdb
go build -o dockerdb ./cmd/dockerdb
```

## Usage

After building the project, you can use the `dockerdb` command followed by the desired action. Here are some examples:

```bash
dockerdb [Choose between MySQL, MariaDB, PostgreSQL, MongoDB, and Redis]
```
After running the command, you wil be taken to a quick setup process. After that, a docker container should pop up if you run:
```bash
docker ps
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.