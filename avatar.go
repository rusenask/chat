package main

import (
	"errors"
)

// ErrNoAvatar is the error when no avatar URL is available
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for client or error if something goes wrong
	// ErrNoAvatarURL is returned if the object is unable to get url for specified
	// client
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar structure, zero init
type AuthAvatar struct{}

// UseAuthAvatar initialises AuthAvatar struct
var UseAuthAvatar AuthAvatar

// GetAvatarURL gets user's avatar or return an error
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}