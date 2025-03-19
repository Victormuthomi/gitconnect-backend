#!/bin/bash

echo "🚀 Setting up GitConnect Backend..."

# Create project folder structure
mkdir -p gitconnect-backend/{config,controllers,middlewares,migrations,models,routes,utils}
cd gitconnect-backend

# Initialize Go module
go mod init gitconnect-backend

# Install dependencies
echo "📦 Installing dependencies..."
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/joho/godotenv

# Create essential files
echo "📂 Creating files..."
touch main.go .env \
      config/config.go \
      controllers/authController.go controllers/profileController.go controllers/postController.go \
      middlewares/authMiddleware.go \
      models/user.go models/profile.go models/post.go models/comment.go \
      routes/authRoutes.go routes/profileRoutes.go routes/postRoutes.go \
      utils/token.go

# Add .env example
cat <<EOT >> .env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=gitconnect
DB_PORT=5432
JWT_SECRET=your_secret_key
EOT

echo "✅ Setup complete! Ready to code. 🚀"

