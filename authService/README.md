# Auth Service

## Overview

The **Auth Service** is a comprehensive authentication and authorization service built with **Go**. It manages user sessions, Multi-Factor Authentication (MFA), OAuth, and user preferences. The service supports secure password management, session tracking, device management, and provides an efficient way to implement security mechanisms such as JWT authentication and OAuth authorization.

Key features include:
- **User Authentication** (Login, Registration)
- **Multi-Factor Authentication (MFA)** via **SMS** or **Email**
- **OAuth Authentication** (OAuth2.0)
- **Session Management** (Tracking and Expiring Sessions with Redis)
- **Device Management** (Tracking devices per user)
- **Security Middleware** (JWT validation, CORS, and error handling)

## Features
- **User Authentication**: Register, login, and authenticate users using JWTs.
- **Multi-Factor Authentication** (MFA): Supports **TOTP** (Time-based One-Time Passwords) and **Backup Codes**.
- **OAuth2 Authentication**: Manage OAuth2 flows for third-party integrations.
- **Session Management with Redis**: Control user sessions with secure session tracking, using Redis for storing session data.
- **Rate Limiting**: Protect endpoints with rate-limiting mechanisms.
- **JWT Authentication**: Secure token-based authentication for stateless APIs.

## Directory Structure

```plaintext
├── auth_grpc                   # gRPC server and services
├── auth_proto_generated        # Generated gRPC files from proto definitions
├── cmd                         # Application entry point
├── config                      # Configuration files (DB, JWT, OAuth, etc.)
├── controllers                 # HTTP controllers for authentication and OAuth endpoints
├── crts                        # SSL certificates for secure communication
├── db                          # Database connections and migrations
├── deployments                 # Kubernetes configurations for deployment
├── docker                      # Docker configurations for local development
├── Dto                         # Data Transfer Objects (DTOs) for request validation
├── go.mod                      # Go module file
├── go.sum                      # Go checksum file
├── internal                    # Core application logic (models, repositories, services)
├── middlewares                 # Custom middlewares (authentication, CORS, error handling)
├── oauth_server                # OAuth server implementation
├── scripts                     # Utility scripts (e.g., migration scripts)
├── test                        # Unit and integration tests
├── utils                       # Utility functions (e.g., JWT signing)
└── README.md                   # This file
