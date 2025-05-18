# WebPaint - Online Drawing App

WebPaint is a web-based drawing application that allows users to create, save, and manage their digital artwork. Built with Go on the backend and modern HTML/CSS/JavaScript on the frontend, it provides a simple yet powerful drawing experience.

## Features

- 🎨 Pixel-perfect drawing canvas with adjustable brush size and color
- 💾 Save drawings with custom titles
- 📂 Gallery to view and manage all saved artwork
- 🗑️ Delete unwanted drawings
- 🔐 User authentication (login/registration)
- 📱 Responsive design that works on desktop devices

## Technology Stack

**Backend:**
- Go (Golang)
- MySQL database
- Standard library HTTP server

**Frontend:**
- Vanilla JavaScript (no frameworks)
- Canvas API for drawing
- Modern CSS with Flexbox and Grid
- Font Awesome icons

## Installation

### Prerequisites
- Go 1.20+
- MySQL 5.7+
- Node.js (for optional frontend development)

### Setup

1. Copy `.env.example` to `.env`
2. Fill in your database credentials
3. Run `go run main.go`
