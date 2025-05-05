package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

const (
	uploadDir  = "./uploads"
	serverPort = ":8081"
)

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	PhotoURL    string `json:"photo_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Birthday    string `json:"birthday"`
	Age         int    `json:"age"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Skype       string `json:"skype"`
}

func main() {
	var err error
	// MySQL connection string
	dsn := "root@tcp(127.0.0.1:3306)/socialmediaapp"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("Error connecting to the database: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		panic("Database ping failed: " + err.Error())
	}
	fmt.Println("✅ Connected to MySQL database")

	r := gin.Default()

	// Setup CORS for frontend interaction
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Public routes
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)
	r.POST("/auth/microsoft", microsoftLoginHandler)
	// Protected routes
	r.POST("/update-profile/:id", updateProfileHandler)
	r.POST("/upload-multiple-photos", uploadMultiplePhotosHandler)
	// Backend route to get user posts/photos
	r.GET("/user-posts/:email", getUserPostsHandler)

	// Serve uploaded images
	r.Static("/uploads", uploadDir)
	r.POST("/delete-photo", deletePhotoHandler)
	r.GET("/api/active-users", getActiveUsersHandler)

	// Start the server
	r.Run(serverPort)
}

// ------------------------------ HANDLERS ------------------------------

func signupHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var exists string
	// Check if the user already exists by email
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&exists)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}
	// Insert user into the database
	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful!"})
}

func loginHandler(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login payload"})
		return
	}

	var user User
	var hashedPassword string
	var nullableFields = struct {
		photoURL, title, description, birthday, address, phone, skype sql.NullString
		age                                                           sql.NullInt64
	}{}

	// Fetch user data from DB by email
	err := db.QueryRow(`SELECT id, name, email, password, photo_url, title, description, birthday, age, address, phone, skype 
		FROM users WHERE email = ?`, credentials.Email).Scan(
		&user.ID, &user.Name, &user.Email, &hashedPassword,
		&nullableFields.photoURL, &nullableFields.title, &nullableFields.description,
		&nullableFields.birthday, &nullableFields.age, &nullableFields.address,
		&nullableFields.phone, &nullableFields.skype)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		}
		return
	}

	// Compare hashed password with the entered password
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Populate user fields with optional data
	user.PhotoURL = safeString(nullableFields.photoURL)
	user.Title = safeString(nullableFields.title)
	user.Description = safeString(nullableFields.description)
	user.Birthday = safeString(nullableFields.birthday)
	user.Age = safeInt(nullableFields.age)
	user.Address = safeString(nullableFields.address)
	user.Phone = safeString(nullableFields.phone)
	user.Skype = safeString(nullableFields.skype)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

func microsoftLoginHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	// Verify the Microsoft token
	claims, err := verifyMicrosoftToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Microsoft token"})
		return
	}

	email, emailOK := claims["preferred_username"].(string)
	name, nameOK := claims["name"].(string)
	if !emailOK || !nameOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claims in token"})
		return
	}

	var user User
	// Check if the user already exists
	err = db.QueryRow("SELECT id, name, email FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email)

	if err == sql.ErrNoRows {
		// Create a new user if they don't exist
		res, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		id, _ := res.LastInsertId()
		user = User{ID: int(id), Name: name, Email: email}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in via Microsoft", "user": user})
}

// ------------------------------ Microsoft Token Verification ------------------------------

func verifyMicrosoftToken(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// NOTE: This does NOT validate the signature. For full prod, use OpenID keys to verify.
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

// ------------------------------ PROFILE UPDATES ------------------------------

func updateProfileHandler(c *gin.Context) {
	userID := c.Param("id")

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form data"})
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	title := c.PostForm("title")
	description := c.PostForm("description")
	birthday := c.PostForm("birthday")
	ageStr := c.PostForm("age")
	address := c.PostForm("address")
	phone := c.PostForm("phone")
	skype := c.PostForm("skype")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		age = 0
	}

	var photoURL string
	file, handler, err := c.Request.FormFile("photo")
	if err == nil {
		os.MkdirAll("./uploads", os.ModePerm)

		extension := filepath.Ext(handler.Filename)
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(handler.Filename))
		filename = strings.TrimSuffix(filename, extension) + ".jpg"
		filepath := filepath.Join("./uploads", filename)

		out, err := os.Create(filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		defer out.Close()

		_, err = out.ReadFrom(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
			return
		}

		photoURL = "/backend/uploads/" + filename
	}

	if photoURL == "" {
		err := db.QueryRow("SELECT photo_url FROM users WHERE id = ?", userID).Scan(&photoURL)
		if err != nil {
			photoURL = ""
		}
	}

	query := `
        UPDATE users
        SET name = ?, email = ?, title = ?, description = ?, birthday = ?, age = ?, address = ?, phone = ?, skype = ?, photo_url = ?
        WHERE id = ?
    `
	_, err = db.Exec(query, name, email, title, description, birthday, age, address, phone, skype, photoURL, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":          userID,
			"name":        name,
			"email":       email,
			"title":       title,
			"description": description,
			"birthday":    birthday,
			"age":         age,
			"address":     address,
			"phone":       phone,
			"skype":       skype,
			"photo_url":   photoURL,
		},
	})
}

// ------------------------------ PHOTO UPLOAD HANDLERS ------------------------------

func uploadMultiplePhotosHandler(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	files := form.File["photos"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	os.MkdirAll(uploadDir, os.ModePerm)

	var uploadedFiles []string

	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filename = strings.TrimSuffix(filename, extension) + ".jpg"
		fullPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		photoURL := "/backend/uploads/" + filename

		_, err := db.Exec("INSERT INTO user_posts (email, photo_url, created_at) VALUES (?, ?, NOW())", email, photoURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert record into database"})
			return
		}

		uploadedFiles = append(uploadedFiles, photoURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Files uploaded successfully",
		"uploaded_files": uploadedFiles,
	})
}

// ------------------------------ USER POSTS ------------------------------

func getUserPostsHandler(c *gin.Context) {
	email := c.Param("email")

	rows, err := db.Query("SELECT photo_url FROM user_posts WHERE email = ?", email)
	if err != nil {
		log.Println("Failed to query database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	defer rows.Close()

	var photos []string
	for rows.Next() {
		var photoURL string
		if err := rows.Scan(&photoURL); err != nil {
			log.Println("Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning posts"})
			return
		}
		photos = append(photos, photoURL)
	}

	if len(photos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No posts found for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

// Delete a photo from user_posts table and the filesystem
func deletePhotoHandler(c *gin.Context) {
	// Extract the email and photo URL from the request
	email := c.DefaultQuery("email", "")
	photoURL := c.DefaultQuery("photo_url", "")
	if email == "" || photoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and photo_url are required"})
		return
	}

	// Remove the file from the filesystem
	filePath := "." + photoURL
	err := os.Remove(filePath)
	if err != nil {
		// If the file doesn't exist or an error occurs while deleting, log it but continue
		log.Printf("Error deleting file: %v\n", err)
	}

	// Remove the photo record from the database
	_, err = db.Exec("DELETE FROM user_posts WHERE email = ? AND photo_url = ?", email, photoURL)
	if err != nil {
		log.Println("Failed to delete photo from database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

func getActiveUsersHandler(c *gin.Context) {
	query := `
	SELECT u.id, u.name, u.photo_url, up.photo_url 
	FROM users u
	JOIN user_posts up ON u.email = up.email
	ORDER BY u.id, up.created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Println("Failed to query active users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active users"})
		return
	}
	defer rows.Close()

	type Post struct {
		PhotoURL string `json:"photo_url"`
	}
	type UserWithPosts struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		PhotoURL string `json:"photo_url"`
		Posts    []Post `json:"posts"`
	}

	userMap := make(map[int]*UserWithPosts)

	for rows.Next() {
		var id int
		var name, userPhotoURL, postPhotoURL string
		if err := rows.Scan(&id, &name, &userPhotoURL, &postPhotoURL); err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		if _, exists := userMap[id]; !exists {
			userMap[id] = &UserWithPosts{
				ID:       id,
				Name:     name,
				PhotoURL: userPhotoURL,
				Posts:    []Post{},
			}
		}
		userMap[id].Posts = append(userMap[id].Posts, Post{PhotoURL: postPhotoURL})
	}

	var users []UserWithPosts
	for _, user := range userMap {
		users = append(users, *user)
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No active users found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Utility Functions for Nullable Fields
func safeString(val sql.NullString) string {
	if val.Valid {
		return val.String
	}
	return ""
}

func safeInt(val sql.NullInt64) int {
	if val.Valid {
		return int(val.Int64)
	}
	return 0
}
