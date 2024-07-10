# Setup Instructions

## Requirements

1. Docker - [Install Docker](https://docs.docker.com/get-docker/)
2. Go - [Install Go](https://golang.org/doc/install)

## Instructions

### 1. Clone the repository

First, clone the repository using Git:

```bash
git clone https://github.com/kshakhnitski/flat-notification-telegram-bot.git
cd flat-notification-telegram-bot
```

### 2. Set Up PostgreSQL Database using Docker-Compose

Ensure Docker is installed on your system. If not, download and install it from Docker's official website.
Once Docker is installed, use Docker-Compose to spin up a PostgreSQL database:

```bash
 docker-compose up -d
```

### 3. Obtain Telegram Bot API Key

1. Open Telegram and search for the `@BotFather` bot.
2. Start a chat with `@BotFather` and use the `/start` command.
3. Use the `/newbot` command to create a new bot.
4. Follow the prompts to name your bot and set a username for it.
5. Once your bot is created, `@BotFather` will provide you with a token. This token is your `TELEGRAM_BOT_API_KEY`.

### 4. Configure environment variables

Create a `.env` file in the root directory of the project and add your environment variables:

```dotenv
DB_PORT=5432
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=flat_bot
TELEGRAM_BOT_API_KEY=your_telegram_bot_api_key
```

Replace `your_telegram_bot_api_key` with the API key you obtained in the previous step.

### 5. Run the project

Finally, start the project by running the following command:

```bash
go run cmd/app/main.go
```

This command will run the Go application, connecting it to the PostgreSQL database and starting your Telegram bot.
