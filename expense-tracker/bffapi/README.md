# BFF API

A Backend For Frontend (BFF) API that provides a unified interface for the bill management web application.

## Description

This service acts as a middleware between the frontend application and backend services (accounts and bill parser). It exposes a simplified API to the frontend while communicating with multiple backend services.

## Features

- Bill management (CRUD operations)
- Receipt parsing with AI
- Unified API for the frontend

## Getting Started

### Prerequisites

- Node.js (v14 or higher)
- npm

### Installation

```bash
npm install
```

### Configuration

Create a `.env` file in the root directory with the following variables:

```env
PORT=3001
ACCOUNTS_API_URL=http://localhost:8080/api/v1
BILL_PARSER_API_URL=http://localhost:8080
```

### Running the API

```bash
npm start
```

For development with auto-reload:

```bash
npm run dev
```

## API Endpoints

### Bill Management

- `GET /api/bills` - Get all bills
- `GET /api/bills/:id` - Get a bill by ID
- `POST /api/bills` - Create a new bill
- `PUT /api/bills/:id` - Update a bill
- `DELETE /api/bills/:id` - Delete a bill

### Bill Parser

- `POST /api/parser/parse` - Parse a receipt image and optionally create a bill

## Backend Services

This BFF API communicates with the following backend services:

- Accounts API: Manages bill data and storage
- Bill Parser API: Processes receipt images using AI
