package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrNoAvatarURL = errors.New("chat: Unable to get avatar URL.")

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

type FileSystemAvatar struct{}

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			m := md5.New()                                                      //Gravatar's guidelines to generate an MD5
			io.WriteString(m, strings.ToLower(useridStr))                       //(after we ensured it was lowercase)
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil //append it to the hardcoded
			// base URL

		}
	}
	return "", ErrNoAvatarURL
}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
