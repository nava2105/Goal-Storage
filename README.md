# nava2105/Goal-Storage

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white)
![GraphQL](https://img.shields.io/badge/GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white)

---

## Table of Contents
1. [General Info](#general-info)
2. [Features](#features)
3. [Technologies](#technologies)
4. [Architecture](#architecture)
5. [Setup & Installation](#setup--installation)
6. [GraphQL API Usage](#graphql-api-usage)
7. [Environment Variables](#environment-variables)
8. [Database Schema](#database-schema)

---

## General Info

This project is a **Go-based backend application** for managing user goals. The service uses **GraphQL** APIs to allow users to perform CRUD (Create, Read, Update, Delete) operations on their goals, which include attributes such as weight, body structure, etc.

The backend is designed to securely integrate:
- **JWT authorization** to validate users through an external authentication service.
- **MongoDB** as the primary database.

It's built with flexibility in mind using the **Abstract Factory Pattern**, enabling seamless extension and modularity.

---

## Features

### Key Functionalities:
1. **User Goal Management**:
    - Create, read, and update personal goals for authenticated users.
    - Goals store user-specific information like weight, and body structure.

2. **Authentication & Authorization**:
    - JWT-based user authentication with tokens validated by an external **Auth API**.
    - Middleware to enforce authorization for all GraphQL queries and mutations.

3. **GraphQL Interface**:
    - Provides a flexible query/mutation-based API which is well-suited for dynamic client requirements.

4. **Highly Modular Design**:
    - **Abstract Factory Pattern** for managing database interaction.
    - Modular components for GraphQL schema, controllers, and repositories.

---

## Technologies

- **Go Programming Language (v1.23)**: Backend code implementation.
- **GraphQL**: API design and interaction.
- **MongoDB**: NoSQL database for scalable data persistence.
- **Gorilla/Mux**: For routing HTTP endpoints.
- **JWT Authentication**: Secure user authorization using access tokens.

---

## Architecture

This project follows a component-based modular architecture:

1. **Controllers**:
    - Manage API requests and responses.
    - Interface with GraphQL and handle business logic.

2. **GraphQL Schema**:
    - Encapsulates all query and mutation configurations.

3. **Factories**:
    - Abstract database interactions and implement the **Abstract Factory Pattern**.
    - Provide flexibility to replace MongoDB with other databases in the future.

4. **Repositories**:
    - Direct interaction with database drivers for CRUD operations.

5. **Middleware**:
    - Enforces security and validates JWT tokens.
    - Attaches authorization headers to the API context.

6. **Initialization**:
    - Handles database connection setup and environment variable loading.

---

## Setup & Installation

### Prerequisites
1. **Go (v1.23 or later)**: Install Go and set up your environment.
2. **MongoDB**: Ensure MongoDB is installed and running.
3. **Environment Configuration**: Required environment variables are defined in a `.env` file (see the [Environment Variables](#environment-variables) section).

---

### Steps to Run
1. **Clone the repository**:
   ```bash
   git clone https://github.com/nava2105/Goal-Storage.git
   cd Goal-Storage
   ```

2. **Download dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:
    - Create a `.env` file in the root directory and configure your app environment (see [Environment Variables](#environment-variables)).

4. **Run the application**:
   ```bash
   go run main.go
   ```

5. **Access the application**:
    - The server runs locally on the default port: `http://localhost:8000`.

---

## GraphQL API Usage

**Note:** All GraphQL queries and mutations require a valid JWT token in the `Authorization` header in the following format:
```
Authorization: Bearer <your_token>
```
The backend provides the following functionality via **GraphQL** queries and mutations:

### Queries

1. **Get User Goals**:
    - Retrieves the goal associated with an authenticated user.
    - Example:
      ```graphql
      query {
        getGoalById(userId: 1) {
          goalId
          userId
          weight
          body_structure
        }
      }

2. **Get UserId**:
    - Example:
      ```graphql
      query {
        userId
      }

---

### Mutations

1. **Create Goal**:
   ```graphql
   mutation {
     createGoal(userId: 1, weight: 70.5, body_structure: "Athletic") {
       goalId
       userId
       weight
       body_structure
     }
   }

2. **Update Goal**:
   ```graphql
   mutation {
     updateGoal(goalId: "1234", userId: 1, weight: 75, body_structure: "Lean") {
       goalId
       userId
       weight
       body_structure
     }
   }

---

## Environment Variables

Define the following settings in the `.env` file or pass them as environment variables:

| Variable    | Description                                | Example                                      |
|-------------|--------------------------------------------|----------------------------------------------|
| `MONGO_URI` | MongoDB connection URI                     | `mongodb://localhost:27017`                  |
| `PORT`      | Port number to run the HTTP server         | `8000`                                       |
| `AUTH_URL`  | External Auth API for user ID verification | `http://localhost:90/api/auth/user_id/token` |

---

## Database Schema

The database consists of the following structure:

### Collection: `ptrainer_goals`

| **Field**        | **Type** | **Description**                    |
|------------------|----------|------------------------------------|
| `_id`            | ObjectId | Unique identifier (auto-generated) |
| `user_id`        | Int64    | User identifier                    |
| `weight`         | Float64  | User's weight (in kg)              |
| `body_structure` | String   | User's body structure description  |

---

## Notes
1. **GraphQL Playground**:
    - You can connect any GraphQL client (e.g., Postman or Insomnia) to test the queries and mutations.

2. **Authentication**:
    - JWT tokens are required for all operations. Attach them as part of the `Authorization` header in the format `Bearer <token>`.

3. **Error Handling**:
    - The API handles malformed requests and provides appropriate HTTP status codes.