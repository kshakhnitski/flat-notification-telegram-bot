# Setup Instructions

1. Clone the repository
   Clone the repository using Git:
   ```bash
   git clone https://github.com/kshakhnitski/flat-notification-telegram-bot.git
   cd flat-notification-telegram-bot
   ```

2. Set up PostgreSQL Database using Docker-Compose
   Make sure you have Docker installed. Then, use Docker-Compose to spin up a PostgreSQL database:
   ```bash
    docker-compose up -d
   ```

3. Configure environment variables
   Create a `.env` file in the root directory of the project and add your environment variables:
    ```dotenv
    DB_PORT=5432
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=password
    DB_NAME=flat_bot
    TELEGRAM_BOT_API_KEY=your_telegram_bot_api_key
    ```

4. Run the project
   Finally, start the project by running the following command:
   ```bash
   go run cmd/app/main.go
   ```
   This command will run the Go application, connecting it to the PostgreSQL database and starting your Telegram bot.
