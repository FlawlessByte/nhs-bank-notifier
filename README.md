# NHS Bank Notifier

NHS Bank Notifier is an application designed to periodically check for available NHS bank shifts, filter shifts for specific units, and send a notification through Telegram. It uses a TTL cache to avoid notifying about shifts that have already been handled.

## Features

- Periodically checks for new shifts from the NHS API.
- Filters shifts based on specific units (e.g., Intensive Care).
- Notifies users via Telegram.
- Customizable via environment variables.
- Docker support for easy deployment.

## Requirements

- Go 1.23 or later
- Docker

## Configuration

The application can be configured using environment variables. Below are the supported variables and their default values:

| Variable              | Default Value                                         | Description                                  |
|-----------------------|-------------------------------------------------------|----------------------------------------------|
| `NHS_USERNAME`        |                                                       | NHS login username                           |
| `NHS_PASSWORD`        |                                                       | NHS login password                           |
| `LOGIN_URL`           | `https://ich.allocate-cloud.co.uk/EmployeeOnlineHealth/ICHLIVE/Login` | NHS login URL                    |
| `TELEGRAM_BOT_TOKEN`  |                                                       | Your Telegram bot API token which you received from BotFather                  |
| `TELEGRAM_CHAT_ID`    |                                                       | The chat ID where notifications will be sent |
| `MAX_TTL`             | `336h`                                                | TTL for cached shifts (default is 2 weeks)   |
| `CHECK_INTERVAL_MINS` | `10`                                                  | Interval between shift checks in minutes     |
| `LOG_LEVEL`           | `WARN`                                                | Log level for the application                |

## Running Locally

1. Clone the repository:
    ```bash
    git clone https://github.com/FlawlessByte/nhs-bank-notifier
    cd nhs-bank-notifier
    ```

2. Copy the `.env.example` file and update it with your own credentials:
    ```bash
    cp .env.example .env
    ```

3. Build and run the application:
    ```bash
    go build -o nhs-bank-notifier ./cmd/nhs-bank-notifier
    ./nhs-bank-notifier
    ```

## Running with Docker

You can build and run the app using Docker:

1. Build the Docker image:
    ```bash
    docker build -t nhs-bank-notifier .
    ```

2. Run the Docker container:
    ```bash
    docker run --env-file .env nhs-bank-notifier
    ```
## Contributions

Contributions are welcome! If you find any bugs, have feature requests, or want to contribute to the project, feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).