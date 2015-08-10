package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	// using FormValue method on http.Request to get user ID that we placed
	// hidden in the HTML form input
	userID := req.FormValue("userid")
	// getting io.Reader type capable of reading uploaded bytes by calling
	// req.FormFile which returns three arguments
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	// file should be passed to methods that require io.Reader since any object
	// that implement multipart.File must also implement io.Reader
	// ioutil.ReadAll method will keep reading from the specified io.Reader until
	// all of the bytes have been received (this is where we receive stream of bytes)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	// using path.Join and path.Ext to build a new filename with userid and
	// file extension from the original filename that we get from
	// multipart.FileHeader
	filename := path.Join("avatars", userID+path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "Successful")
}
