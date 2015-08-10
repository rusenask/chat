package main

// ErrNoAvatar is the error when no avatar URL is available
// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for client or error if something goes wrong
	// ErrNoAvatarURL is returned if the object is unable to get url for specified
	// client
	GetAvatarURL(c *client) (string, error)
}
