package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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

// GravatarAvatar structure
type GravatarAvatar struct{}

// UseGravatar variable for storing GravatarAvatar
var UseGravatar GravatarAvatar

// GetAvatarURL implements fetching avatar image
func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			// creating new md5 hasher from crypto package, implements io.Writer interface
			// TODO: hashing email every time is not very efficient
			m := md5.New()
			// ensuring email is lower case and generate md5 hash
			// writing a string of bytes to hasher through io.WriteString
			io.WriteString(m, strings.ToLower(emailStr))
			// caling Sum on hasher returns the current hash for the bytes written
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
