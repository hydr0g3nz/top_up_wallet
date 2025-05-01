## Getting Started

Follow these steps to get the project up and running locally using Docker Compose.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/hydr0g3nz/top_up_wallet.git
    cd top_up_wallet
    ```
    Replace `<repository_url>` with the actual URL of your repository and `<repository_directory>` with the resulting directory name.

2.  **Create the Environment File:**

    Create a file named `.env` in the root directory of the project (where `docker-compose.yml` is located). Copy the following content into this file.

    **Important:** Change the default values for `DB_PASSWORD` and `REDIS_PASSWORD` for any environment other than local development.

    ```env
    # Server settings
    PORT=8080
    SERVER_READ_TIMEOUT=10
    SERVER_WRITE_TIMEOUT=10
    SERVER_HOST=localhost

    # Database settings
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=user
    DB_PASSWORD=pass
    DB_NAME=topup_wallet
    DB_SSLMODE=disable

    # Cache settings
    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=pass
    REDIS_DB=0

    # Logging
    LOG_LEVEL=info


    ```


3.  **Build and Run with Docker Compose:**

    Open your terminal in the project's root directory and run the following command:

    ```bash
    docker compose up -d
    ```

    * `docker compose up`: Starts the services defined in `docker-compose.yml`.
    * `-d`: Runs the containers in detached mode (in the background).

    This command will:
    * Build the Docker images for the services (if not already built).
    * Start the containers defined in `docker-compose.yml`.
    * Run the containers in detached mode (in the background).

4.  **Verify Containers are Running:**

    You can check the status of the containers using:

    ```bash
    docker compose ps
    ```

    * `docker compose ps`: Lists the status of the services defined in `docker-compose.yml`.

5.  **Access the API:**

    Once the containers are up and running, you can access the API at the following URL: `http://localhost:8080/api/v1`.


## Configuration

The project's behavior is primarily controlled by environment variables, configured via the `.env` file when using Docker Compose.

* `PORT`: Internal port the Go application listens on (exposed via Docker).
* `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`: Database connection details.
* `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`: Redis cache connection details.
* `LOG_LEVEL`: Log level for the application.

## Stopping the Project

To stop the running containers:

```bash
docker compose down
