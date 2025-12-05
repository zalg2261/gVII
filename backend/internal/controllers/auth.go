package controllers

import (
    "os"
    "strconv"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "github.com/zalg2261/bioskop/backend/internal/db"
    "github.com/zalg2261/bioskop/backend/internal/models"
)

func getJWTKey() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return []byte("secret123") // default for development
    }
    return []byte(secret)
}

func Register(c *fiber.Ctx) error {
    input := new(models.User)
    if err := c.BodyParser(input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
    }
    hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
    input.Password = string(hash)
    if err := db.DB.Create(&input).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "registered"})
}

func Login(c *fiber.Ctx) error {
    input := new(models.User)
    if err := c.BodyParser(input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
    }
    var user models.User
    if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "user not found"})
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "wrong password"})
    }
    claims := jwt.MapClaims{
        "id":  strconv.Itoa(int(user.ID)),
        "role": user.Role,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, _ := token.SignedString(getJWTKey())
    return c.JSON(fiber.Map{
        "token": t,
        "user": fiber.Map{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
            "role":  user.Role,
        },
    })
}
