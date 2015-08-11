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
	GetAvatarURL(u ChatUser) (string, error)
}

// AuthAvatar structure, zero init
type AuthAvatar struct{}

// UseAuthAvatar initialises AuthAvatar struct
var UseAuthAvatar AuthAvatar

// GetAvatarURL gets user's avatar or return an error
func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) > 0 {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar structure
type GravatarAvatar struct{}

// UseGravatar variable for storing GravatarAvatar
var UseGravatar GravatarAvatar

// GetAvatarURL implements fetching avatar image
func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar structure
type FileSystemAvatar struct{}

// UseFileSystemAvatar helper variable
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL implements using avatar image from uploaded user files
func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		// checking if file is directory - if yes, moving on
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			// if file matches userid - returning it
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}
