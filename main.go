// TTW Software Team
// Mathis Van Eetvelde
// 2021-present

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	scope = "https://www.googleapis.com/auth/drive.file"
	filenameInput = "filename"
	nameInput = "name"
	folderIdInput = "folderId"
	credentialsInput = "credentials"
)

func main() {

	// get filename argument from action input
	filename := githubactions.GetInput(filenameInput)
	if filename == "" {
		missingInput(filenameInput)
	}

	// get name argument from action input
	name := githubactions.GetInput(nameInput)

	// get folderId argument from action input
	folderId := githubactions.GetInput(folderIdInput)
	if folderId == "" {
		missingInput(folderIdInput)
	}

	// get base64 encoded credentials argument from action input
	credentials := githubactions.GetInput(credentialsInput)
	if credentials == "" {
		missingInput(credentialsInput)
	}
	// add base64 encoded credentials argument to mask
	githubactions.AddMask(credentials)

	// decode credentials to []byte
	decodedCredentials, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("base64 decoding of 'credentials' failed with error: %v", err))
	}

	creds := strings.TrimSuffix(string(decodedCredentials), "\n")

	// add decoded credentials argument to mask
	githubactions.AddMask(creds)

	// fetching a JWT config with credentials and the right scope
	conf, err := google.JWTConfigFromJSON([]byte(creds), scope)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("fetching JWT credentials failed with error: %v", err))
	}

	// instantiating a new drive service
	ctx := context.Background()
	svc, err := drive.New(conf.Client(ctx))
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("opening file with filename: %v failed with error: %v", filename, err))
	}

	// decide name of file in GDrive
	if name == "" {
		name = file.Name()
	}

	f := &drive.File{
		Name:    name,
		Parents: []string{folderId},
	}
	githubactions.Warningf(fmt.Sprintf("File to upload: %v", f))
	uploaded, err := svc.Files.Create(f).Media(file).Do()
	githubactions.Warningf(fmt.Sprintf("File uploaded: %v", uploaded))
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("creating file: %+v failed with error: %v", f, err))
	}
	githubactions.Warningf(fmt.Sprintf("f.WebViewLink: %v", f.WebViewLink))
	githubactions.Warningf(fmt.Sprintf("f.WebContentLink: %v", f.WebContentLink))
	githubactions.Warningf(fmt.Sprintf("f.ExportLinks: %v", f.ExportLinks))

	githubactions.Warningf(fmt.Sprintf("uploaded file: %v", uploaded))
	if uploaded != nil {

		githubactions.Warningf(fmt.Sprintf("uploaded.WebViewLink: %v", uploaded.WebViewLink))
		githubactions.Warningf(fmt.Sprintf("uploaded.WebContentLink: %v", uploaded.WebContentLink))
		githubactions.Warningf(fmt.Sprintf("uploaded.ExportLinks: %v", uploaded.ExportLinks))
		link := fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=drivesdk", uploaded.Id)
		githubactions.SetOutput("link", link)
		githubactions.Warningf(fmt.Sprintf("link: %v", link))

	} else {
		githubactions.Errorf(fmt.Sprintf("No error uploading file but no file uploaded (?): %v - %v", f, uploaded))
	}

}

func missingInput(inputName string) {
	githubactions.Fatalf(fmt.Sprintf("missing input '%v'", inputName))
}
