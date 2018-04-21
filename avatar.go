package main

import (
	"errors"
	"io/ioutil"
	"path"
)

var ErrNoAvatarURL = errors.New("chat: Unable to get avatar URL.")

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

var UseAuthAvatar AuthAvatar //nil atm,  later assign the UseAuthAvatar
// variable to any field looking for an avatar interface type

type FileSystemAvatar struct{}

type Avatar interface {
	GetAvatarURL(ChatUser) (string, error)
}
type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
} //TDD

// func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
// 	if userid, ok := c.userData["userid"]; ok {
// 		if useridStr, ok := userid.(string); ok {
// 			m := md5.New()                                                      //Gravatar's guidelines to generate an MD5
// 			io.WriteString(m, strings.ToLower(useridStr))                       //(after we ensured it was lowercase)
// 			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil //append it to the hardcoded
// 			// base URL

// 		}
// 	}
// 	return "", ErrNoAvatarURL
// }

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) > 0 {
		return url, nil
	}
	return "", ErrNoAvatarURL
}
