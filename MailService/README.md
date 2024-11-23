# MailService

A service to send emails using Nodemailer and TypeScript, integrated with Kafka for receiving messages from a broker and sending emails based on those messages.

## Features

- Send emails via SMTP using [Nodemailer](https://www.nodemailer.com/)
- Kafka integration for receiving messages from a broker
- Written in TypeScript
- Supports dynamic email content based on received messages
- Configurable SMTP and Kafka settings
