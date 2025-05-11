# Snippetbox

Snippetbox is a web application for managing code snippets, built using Go and based on the tutorial from _Let's Go_ by Alex Edwards. The project demonstrates various web development concepts, including user authentication, form validation, secure password handling, and CSRF protection.

This application allows users to sign up, log in, and manage their snippets, with a focus on secure and scalable design.

## Features

- **User Authentication**: Users can sign up, log in, and access their snippets.
- **Form Validation**: Ensures that user input is correct before processing (e.g., email format, password strength).
- **Error Handling**: Comprehensive error handling with detailed feedback for the user.
- **CSRF Protection**: Prevents cross-site request forgery attacks.
- **Session Management**: Persistent login sessions with secure cookies.
- **Secure Password Storage**: Passwords are stored securely using bcrypt hashing.
- **Database Interaction**: Uses MySQL for persistent storage of users and snippets.
- **HTTP Redirects**: Handles HTTP redirects with proper status codes (e.g., `303 See Other` after form submissions).
- **Flash Messages**: Provides feedback to the user after form submissions.

## Installation

To run the Snippetbox project locally, follow these steps:

### 1. Clone the Repository

```bash
git clone https://github.com/ashab-k/snippetbox.git
cd snippetbox
```

### 2. Run using docker

```bash
docker compose up --build
```
