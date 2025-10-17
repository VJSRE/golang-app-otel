package controllers

import (
	"net/http"

	"github.com/SigNoz/sample-golang-app/models"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// GET /books
// Find all books
func FindBooks(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Fetching all books")
	var books []models.Book
	models.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Fetching single book")
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		span.AddEvent("Book not found", trace.WithAttributes(
			attribute.String("id", c.Param("id")),
		))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Creating book")
	// Validate input
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		span.AddEvent("Book creation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := models.Book{Title: input.Title, Author: input.Author}
	models.DB.Create(&book)
	span.AddEvent("Book created", trace.WithAttributes(
		attribute.Int("book_id", int(book.ID)),
	))
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PATCH /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Updating book")
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		span.AddEvent("Book not found", trace.WithAttributes(
			attribute.String("id", c.Param("id")),
		))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		span.AddEvent("Book update failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&book).Updates(input)
	span.AddEvent("Book updated successfully", trace.WithAttributes(
		attribute.Int("book_id", int(book.ID)),
	))
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "books"))
	span.AddEvent("Deleting book")
	// Get model if exist
	var book models.Book
	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		span.AddEvent("Book not found", trace.WithAttributes(
			attribute.String("id", c.Param("id")),
		))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&book)
	span.AddEvent("Book deleted", trace.WithAttributes(attribute.Int("book_id", int(book.ID))))
	c.JSON(http.StatusOK, gin.H{"data": true})
}
