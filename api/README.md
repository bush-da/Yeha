# Yeha API

Welcome to the Yeha API repository! This API serves as the backend for the Yeha blogging platform, enabling users to interact with the platform through various endpoints.

## Features
- User Management: Registration, login, profile updates, and deletion.
- Blogging: Create, update, delete, and retrieve posts.
- Comments: Add, delete, and fetch comments under posts.
- Likes: Like and unlike posts.
- Admin Controls: Manage reported posts, delete users, and oversee content.
- JWT Authentication: Secure access to protected routes.

## Technologies Used
- **Programming Language:** Go
- **Framework:** Gin
- **Database:** MySQL
- **ORM:** GORM
- **Authentication:** JWT
- **Unique Identifiers:** UUID

## Prerequisites
To run this API, ensure the following are installed on your system:
- Go (version 1.20+)
- MySQL (version 8.0+)
- Postman (optional, for testing API endpoints)

🛠️ ## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/bush-da/Yeha/
   cd api
   ```
2. Set up environment variables in a `.env` file:
   ```
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=yeha
   DB_HOST=localhost
   DB_PORT=3306
   JWT_SECRET=your_jwt_secret
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run database migrations (if applicable).
5. Start the server:
   ```bash
   go run main.go
   ```

## Folder Structure
```
api/
├── controllers/    # Business logic and request handling
├── database/       # Database connection setup
├── middleware/     # Authentication and request validation
├── migrate/        # migrate the database
├── tests/          # testing api scripts
├── models/         # Database schemas and relationships
├── routes/         # Route definitions and grouping
├── main.go         # Entry point of the application
├── go.mod          # Dependency management
├── go.sum          # Dependency locks
```

## API Endpoints
### Authentication
- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - Login and obtain a token

### User Management
- `GET /api/users` - Get all users (Admin only)
- `GET /api/users/:id` - Get a user by ID
- `PUT /api/users/:id` - Update user details
- `DELETE /api/users/:id` - Delete a user

### Posts
- `POST /api/posts` - Create a new post
- `GET /api/posts` - Retrieve all posts
- `GET /api/posts/:id` - Retrieve a post by ID
- `PUT /api/posts/:id` - Update a post
- `DELETE /api/posts/:id` - Delete a post

### Comments
- `POST /api/posts/:id/comments` - Add a comment to a post
- `GET /api/posts/:id/comments` - Get comments for a post
- `DELETE /api/comments/:id` - Delete a comment

### Likes
- `POST /api/posts/:id/like` - Like a post
- `DELETE /api/posts/:id/like` - Unlike a post

### Admin Controls
- `GET /api/reports` - View reported posts
- `PUT /api/reports/:id` - Mark a report as reviewed
- `DELETE /api/posts/:id` - Delete a reported post
- `DELETE /api/users/:id` - Delete a user

## Development
1. Run the server locally:
   ```bash
   go run main.go
   ```
2. Use tools like Postman or Curl to test endpoints.
3. Follow RESTful practices for extending or adding new features.

## Testing
- Use the provided Postman collection for testing the API endpoints.
- Run unit tests:
   ```bash
   go test ./...
   ```

🤝## Contributing
We welcome contributions! Please follow the steps below:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-name`).
3. Make your changes and commit (`git commit -m 'Add feature-name'`).
4. Push to the branch (`git push origin feature-name`).
5. Open a Pull Request.

## License
This project is licensed under the MIT License.

---

Feel free to reach out with any questions or suggestions. Happy coding!
