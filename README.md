# Gingo Helpers

A comprehensive Go library that provides helper utilities for building web applications with the Gin framework. This library includes validation middleware, file handling, JWT authentication, database utilities, and testing helpers.

## Features

- **Request Validation**: Automatic validation middleware for JSON and multipart form data
- **File Handling**: Easy file upload and management utilities
- **JWT Authentication**: Token generation and validation
- **Database Utilities**: GORM integration with testing support
- **Test Helpers**: Comprehensive testing utilities with Docker support
- **Type Conversion**: Automatic conversion between validation structs and database models

## Installation

```bash
go get github.com/sajad-dev/gingo-helpers
```

## Quick Start

### 1. Bootstrap Configuration

```go
package main

import (
    "github.com/sajad-dev/gingo-helpers/core/bootstrap"
    "github.com/sajad-dev/gingo-helpers/types"
)

func main() {
    bootstrap.Boot(types.Bootsterap{
        Config: types.ConfigUtils{
            STORAGE_PATH: "./storage",
            JWT:          "your-jwt-secret",
            DATABASE:     []any{&User{}, &Post{}}, // Your models
        },
    })
}
```

### 2. Request Validation

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/sajad-dev/gingo-helpers/core/validation"
)

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func createUserHandler(ctx *gin.Context) {
    req := ctx.MustGet("validated_data").(CreateUserRequest)
    // Your business logic here
    ctx.JSON(200, gin.H{"message": "User created successfully"})
}

func ValidationMiddleware(validationStruct any) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        validated, err := validation.SwitchHeader(ctx, validationStruct)
        if err != nil {
            ctx.JSON(400, gin.H{"error": err.Error()})
            return
        }
        
        ctx.Set("validated_data", validated.Interface())
        ctx.Next()
    }
}

func main() {
    r := gin.Default()
    
    r.POST("/users", 
        ValidationMiddleware(CreateUserRequest{}),
        createUserHandler,
    )
    
    r.Run(":8080")
}
```

### 3. File Upload Handling

```go
package main

import (
    "mime/multipart"
    "github.com/gin-gonic/gin"
    "github.com/sajad-dev/gingo-helpers/utils"
)

type FileUploadRequest struct {
    Title string                  `json:"title"`
    Files []*multipart.FileHeader `json:"files"`
}

func uploadHandler(ctx *gin.Context) {
    req := ctx.MustGet("validated_data").(FileUploadRequest)
    
    // Save files
    savedPaths, err := utils.SaveFile(ctx, "uploads", req.Files, "document")
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(200, gin.H{
        "message": "Files uploaded successfully",
        "paths": savedPaths,
    })
}
```

### 4. JWT Authentication

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/sajad-dev/gingo-helpers/utils"
)

func loginHandler(ctx *gin.Context) {
    // Authenticate user logic here
    userID := 123
    email := "user@example.com"
    
    // Create JWT token
    token, err := utils.CreateJWT(map[string]any{
        "user_id": userID,
        "email":   email,
    })
    
    if err != nil {
        ctx.JSON(500, gin.H{"error": "Failed to create token"})
        return
    }
    
    ctx.JSON(200, gin.H{"token": token})
}

func AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        tokenString := ctx.GetHeader("Authorization")
        if tokenString == "" {
            ctx.JSON(401, gin.H{"error": "Authorization header required"})
            ctx.Abort()
            return
        }
        
        // Remove "Bearer " prefix if present
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
        
        claims, valid, err := utils.ValidJWT(tokenString)
        if err != nil || !valid {
            ctx.JSON(401, gin.H{"error": "Invalid token"})
            ctx.Abort()
            return
        }
        
        ctx.Set("user_claims", claims.Parameters)
        ctx.Next()
    }
}
```

### 5. Database Model Conversion

```go
package main

import (
    "github.com/sajad-dev/gingo-helpers/utils"
    "gorm.io/datatypes"
)

type UserRequest struct {
    Name        string      `json:"name"`
    Preferences []Preference `json:"preferences"`
}

type OtherData struct {
    Avatar string `json:"avatar"`
}

type User struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name"`
    Preferences datatypes.JSON `json:"preferences"`
    Avatar      string         `json:"avatar"`
}

func createUser(userReq *UserRequest, otherData *OtherData) error {
    var user User
    
    // Convert validation struct to database model
    err := utils.ConvertValidationToTable(userReq, otherData, &user)
    if err != nil {
        return err
    }
    
    // Save to database
    // db.Create(&user)
    
    return nil
}
```

## API Reference

### Core Bootstrap

- `bootstrap.Boot(types.Bootsterap)`: Initialize the library with configuration

### Validation

- `validation.SwitchHeader(ctx, validationStruct)`: Validate request based on Content-Type
- Supports both JSON and multipart form data

### File Utilities

- `utils.SaveFile(ctx, namespace, files, name)`: Save uploaded files
- `utils.SetMultipartFields(ctx, obj)`: Set multipart form fields

### JWT Utilities

- `utils.CreateJWT(claims)`: Create JWT token
- `utils.ValidJWT(token)`: Validate JWT token
- `utils.PasswordHashing(password)`: Hash passwords using SHA-256

### Database Utilities

- `utils.ConvertValidationToTable(validation, other, table)`: Convert structs to database models
- `utils.SetupTestDB()`: Setup SQLite in-memory database for testing
- `utils.SetupDB()`: Setup MySQL database with Docker for testing

### Testing Utilities

- `utils.SendRequest(Request)`: Send HTTP requests for testing
- `utils.CreateServer(path, method, handlers, port)`: Create test server
- `utils.CheckValidationErr(response, field, tag)`: Check validation errors in responses

### Other Utilities

- `utils.GenerateToken()`: Generate random tokens
- `utils.ConvertStringToKind(str, kind)`: Convert strings to specific Go types

## Testing

The library includes comprehensive testing utilities:

```go
package main

import (
    "testing"
    "net/http"
    "github.com/sajad-dev/gingo-helpers/utils"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    server := utils.CreateServer("/users", "post", []gin.HandlerFunc{
        ValidationMiddleware(CreateUserRequest{}),
        createUserHandler,
    }, 8080)
    
    response, err := utils.SendRequest(utils.Request{
        Method: http.MethodPost,
        Path:   "/users",
        Engin:  server,
        Headers: map[string]string{"Content-Type": "application/json"},
        Inputs: CreateUserRequest{
            Name:     "John Doe",
            Email:    "john@example.com",
            Password: "password123",
        },
    })
    
    assert.NoError(t, err)
    assert.Equal(t, 200, response.Code)
}
```

## Configuration

The library uses a centralized configuration system:

```go
type ConfigUtils struct {
    STORAGE_PATH string   // Path for file storage
    JWT          string   // JWT secret key
    IMAGE_TEST   string   // Test image path (auto-generated)
    PROJECT_PATH string   // Project root path (auto-detected)
    DATABASE     []any    // Database models for migration
}
```

## Dependencies

- **Gin**: Web framework
- **GORM**: ORM library
- **JWT-Go**: JWT implementation
- **Dockertest**: Docker testing utilities
- **Testify**: Testing toolkit

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For support and questions, please open an issue in the GitHub repository.
