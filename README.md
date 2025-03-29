# DockerDB

DockerDB is a command-line utility designed to simplify the setup and management of Docker containers for various databases. This tool provides an easy way to create, start, and manage containers for popular databases such as MySQL, MariaDB, PostgreSQL, MongoDB, and Redis.

## Features

- **Easy Setup**: Quickly set up Docker containers for different databases with simple commands.
- **Multiple Database Support**: Supports MySQL, MariaDB, PostgreSQL, MongoDB, and Redis.
- **Docker API Integration**: Interact with the Docker API to manage containers seamlessly.

## Support

### DockerDB supports:
- MySQL
- MariaDB
- ProgreSQL
- MongoDB
- Redis

## Installation
To install DockerDB, clone the repository and build the project:

```bash
git clone https://github.com/Tygo-lex/dockerdb.git
cd dockerdb
chmod +x install.sh
./install.sh
```
Note that you need Go and Docker installed for this to work!
## Usage

After building the project, you can use the `dockerdb` command followed by the desired action. Here is a examples:

```bash
dockerdb mysql
```
After running the command, you wil be taken to a quick setup process. After that, a docker container should pop up if you run:
```bash
docker ps
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.
