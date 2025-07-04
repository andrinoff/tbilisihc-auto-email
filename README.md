# Tbilisi HC Auto Email

[![Go Version](https://img.shields.io/badge/go-1.23.1+-blue.svg)](https://golang.org/dl/)
[![License: Unlicense](https://img.shields.io/badge/License-Unlicense-gray.svg)](https://opensource.org/licenses/unlicense)
[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Ftbilisihc-auto-email%2Ftbilisihc-auto-email)

A simple, Vercel-hosted serverless function to send emails automatically. This project is perfect for handling automated email tasks like sending welcome emails, notifications, or other transactional messages for the Tbilisi Hack Club.

---

## Features

-   **Serverless:** Deployed on Vercel for easy scaling and management.
-   **Go-powered:** Written in Go for performance and reliability.
-   **HTML Emails:** Sends nicely formatted HTML emails using a template.
-   **CORS Support:** Configured to accept requests from specific domains, making it a secure backend for your frontend applications.

---

## API

The service exposes a single endpoint that accepts POST requests to send an email.

-   **Endpoint:** `/api/welcome`
-   **Method:** `POST`
-   **Content-Type:** `application/json`

### Request Body

| Field       | Type   | Description                                |
| :---------- | :----- | :----------------------------------------- |
| `recipient` | string | The email address of the recipient.        |
| `subject`   | string | The subject of the email.                  |
| `message`   | string | The HTML content of the email body.        |

### Example Request

```bash
curl -X POST \
  [https://your-vercel-deployment-url.vercel.app/api/welcome](https://your-vercel-deployment-url.vercel.app/api/welcome) \
  -H 'Content-Type: application/json' \
  -d '{
    "recipient": "test@example.com",
    "subject": "Welcome to the Tbilisi Hack Club!",
    "message": "<h1>Hello!</h1><p>This is a test message to welcome you to our community.</p>"
  }'
```

---

## Environment Variables

To run this project, you will need to add the following environment variables to your Vercel project:

-   `YAHOO_EMAIL`: Your Yahoo email address used for sending emails.
-   `YAHOO_APP_PASSWORD`: Your Yahoo app password for authentication.

This uses particularly YAHOO, but you can use anything else.

---

## Deployment

You can deploy this project to your own Vercel account with a single click:

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Ftbilisihc-auto-email%2Ftbilisihc-auto-email)

---

## Local Development

To run this project locally, you'll need to have Go and the Vercel CLI installed.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/tbilisihc-auto-email/tbilisihc-auto-email.git
    cd tbilisihc-auto-email
    ```

2.  **Install Vercel CLI:**
    ```bash
    npm install -g vercel
    ```

3.  **Create a `.env` file** in the root directory and add your environment variables:
    ```
    YAHOO_EMAIL="your-email@yahoo.com"
    YAHOO_APP_PASSWORD="your-yahoo-app-password"
    ```

4.  **Start the development server:**
    ```bash
    vercel dev
    ```

Your serverless function will now be running on a local port, which will be displayed in your terminal.

---

## Contributing

Contributions are welcome! If you have any ideas, suggestions, or bug reports, please feel free to open an issue or submit a pull request.


---

## License

This project is licensed under the Unlicense License - see the [LICENSE](LICENSE) file for details.
