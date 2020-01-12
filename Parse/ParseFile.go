package Parse

import (
	"../Misc"
	"bufio"
	"errors"
	"os"
)

//Parsing -p and -P supplied parameters
func ParsePass(h *Misc.HostInfo) ([]string, error) {
	switch {
	case h.Passfile != "" && h.Password != "":
		return nil, errors.New("-p and -P cannot exist at the same time")
	case h.Passfile == "" && h.Password == "":
		return nil, errors.New("Please use -p or -P to specify the username")
	case h.Passfile == "" && h.Password != "":
		return []string{h.Password}, nil
	case h.Passfile != "" && h.Password == "":
		file, err := Readfile(h.Passfile)
		return file, err
	}
	return nil, errors.New("unknown error")
}

//Parsing -u and -U supplied parameters
func ParseUser(h *Misc.HostInfo) ([]string, error) {
	switch {
	case h.Userfile != "" && h.Username != "":
		return nil, errors.New("-u and -U cannot exist at the same time")
	case h.Userfile == "" && h.Username == "":
		return nil, errors.New("Please use -u or -U to specify the username")
	case h.Userfile == "" && h.Username != "":
		return []string{h.Username}, nil
	case h.Userfile != "" && h.Username == "":
		file, err := Readfile(h.Userfile)
		return file, err
	}
	return nil, errors.New("unknown error")
}

//Read the contents of the file
func Readfile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("There was an error opening the file")
	}
	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	return content, nil
}
