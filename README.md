# Scheduler API

```
scheduler-api
├── api
│   ├── controllers
│   │   ├── auth.go  # Handles authentication requests.
│   │   ├── event.go  # Handles event requests.
│   │   ├── user.go  # Handles user requests.
│   │   └── user_event.go  # Handles user event requests.
│   ├── models
│   │   ├── event.go  # Represents an event.
│   │   └── user.go  # Represents a user.
│   ├── middleware
│   │   └── auth.go  # Authenticates requests.
│   └── routes.go  # Maps HTTP requests to controllers.
├── initializers
│   ├── db.go  # initializes the database connection.
│   └── logger.go  # initializes the logger.
├── main.go
├── go.mod
└── go.sum
└── test
    └── api
        ├── controllers
        │   ├── auth_test.go  # Tests the auth controller.
        │   ├── event_test.go  # Tests the event controller.
        │   └── user_test.go  # Tests the user controller.
        └── models
            ├── event_test.go  # Tests the event model.
            └── user_test.go  # Tests the user model.
```