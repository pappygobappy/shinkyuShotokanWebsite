package middleware

import (
	"fmt"
	"log"
	"os"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *fiber.Ctx) error {
	//Get the cookie off request
	tokenString := c.Cookies("Authorization")

	//Decode and validate

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("HMAC_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			log.Print("token expired")
			return c.Redirect("/")
		}

		//Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			log.Print("no user found")
			return c.Redirect("/")
		}
		//Attach to req
		c.Locals("user", user)

		//Continue
		return c.Next()
	} else {
		fmt.Println(err)
		return c.Redirect("/")
	}
}

func AttachUser(c *fiber.Ctx) error {
	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}
	c.Locals("hxRequest", hxRequest)

	//Get the cookie off request
	tokenString := c.Cookies("Authorization", "NoAuth")

	if tokenString == "NoAuth" {
		return c.Next()
	}

	//Decode and validate

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("HMAC_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			log.Print("token expired")
			return c.Next()
		}

		//Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			log.Print("no user found")
			return c.Next()
		}
		//Attach to req
		c.Locals("user", user)

		//Continue
		return c.Next()
	} else {
		fmt.Println(err)
		return c.Next()
	}
}
