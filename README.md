# Maestro - The Jazz CLI Tool

Maestro is the command-line tool designed for the Jazz framework, inspired by Laravel's Artisan. It helps automate development tasks, making it easier to manage projects, generate code, and run common operations. Maestro enhances productivity by simplifying repetitive tasks and providing a smooth workflow.

## Features

- **Project Scaffolding**: Create new Jazz projects quickly with pre-configured structures.
- **Code Generation**: Generate controllers, models, and other necessary files with ease.
- **Database Migrations**: Manage database migrations, including creating, running, and rolling back migrations.
- **Development Server**: Easily start a local development server for rapid iteration.

## Getting Started

### Prerequisites

- **Debian-based system**
- **Go installed**

### Installation

To install Maestro on a Debian-based system, you can use the provided installation script:

1. Clone this repository or download the script.
2. Make the script executable and run it:

   ```bash
   chmod +x maestro_install_script.sh
   ./maestro_install_script.sh
   ```

This script will:

- Install Go if it is not already installed.
- Clone the Maestro repository.
- Build the Maestro CLI.
- Move the binary to `~/.local/share/maestro`.
- Add `MAESTRO_HOME` and update `PATH` in `~/.bashrc`.

### Usage

After installing Maestro, you can use it by running:

```bash
maestro <command>
```

Available commands:

- **new**: Create a new project.

  ```bash
  maestro new project-name
  ```

  This command will set up a new Jazz project with the recommended directory structure and dependencies.

- **migrate**: Run database migrations.

  ```bash
  maestro migrate
  ```

  This command will run the migration scripts located in the `backend/migrations` directory.

- **serve**: Start the development server.
  ```bash
  maestro serve [host] [port]
  ```
  By default, it starts on `localhost:8080`. You can optionally specify the host and port, for example:
  ```bash
  maestro serve 0.0.0.0 9090
  ```

## Documentation

Detailed documentation will be provided to guide you through using each of Maestro's commands effectively. Stay tuned!

## Contributing

We welcome contributions to make Maestro even better. Whether it's bug reports, feature suggestions, or pull requests, all feedback is highly appreciated.

## License

Maestro is open-source software available under the MIT License.
