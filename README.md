# Notify

## Description

Notify is an application designed to facilitate notifications for GitHub repositories. By creating webhooks within GitHub, you can choose which notifications you would like to receive in your email. Notify monitors your repository and, for each desired change, sends you an email detailing what happened in your repository.

## Features

- **Webhook Configuration:** Set up webhooks in GitHub repositories to receive notifications.
- **Notification Selection:** Choose which types of notifications you want to receive (commits, pull requests, issues, etc.).
- **Email Delivery:** Receive emails about specific changes in your repositories.
- **Continuous Monitoring:** The backend continuously monitors the configured repositories and triggers notifications as needed.

## Technologies Used

This repository is the backend of the Notify application and uses the following stack:

- **Golang:** Primary language used for application development.
- **Gorilla Mux:** An HTTP request router and dispatcher for Golang.
- **Goth:** A package for authentication with multiple providers, such as GitHub.
- **TursoDB:** The database used for storage.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
