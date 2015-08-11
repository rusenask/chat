package main

import (
	"errors"
	"io/ioutil"
	"path"
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
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// FileSystemAvatar structure
type FileSystemAvatar struct{}

// UseFileSystemAvatar helper variable
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL implements using avatar image from uploaded user files
func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			// getting all files/directories in specified directorie
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				// checking if file is directory - if yes, moving on
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					// if file matches userid - returning it
					if match, _ := path.Match(useridStr+"*", file.Name()); match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
