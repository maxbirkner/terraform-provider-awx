package awx

import (
	"log"
	"os"
	"testing"
)

type TestRow struct {
	data   map[string]interface{}
	params map[string]string
}

var (
	awxHostname string
	awxUsername string
	awxPassword string

	awxClient *AWX

	credentialsServiceTestTable = []*TestRow{{
		data: map[string]interface{}{
			"credential_type": 5,
			"inputs": map[string]interface{}{
				"username": "badusername",
				"password": "badpassword",
			},
			"name":         "credential_01",
			"organization": 1,
		},
		params: map[string]string{},
	},
	}
)

func TestMain(m *testing.M) {
	var err error
	awxHostname = os.Getenv("GOAWX_HOSTNAME")
	awxUsername = os.Getenv("GOAWX_USERNAME")
	awxPassword = os.Getenv("GOAWX_PASSWORD")

	if awxHostname == "" {
		log.Fatal("no AWX hostname provided")
	}

	if awxUsername == "" {
		log.Fatal("no AWX username provided")
	}

	if awxPassword == "" {
		log.Fatal("no AWX password provided")
	}

	awxClient, err = NewAWX(awxHostname, awxUsername, awxPassword, nil)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestCredentialsService(t *testing.T) {
	var createResponse *Credential

	for _, tt := range credentialsServiceTestTable {
		t.Run("Create", func(t *testing.T) {
			var err error
			createResponse, err = awxClient.CredentialsService.CreateCredentials(tt.data, tt.params)
			if err != nil {
				t.Error(err)
			}

			if createResponse.Name != tt.data["name"] {
				t.Errorf("Expecting %s but got %s", tt.data["name"], createResponse.Name)
			}
		})

		t.Run("Fetch", func(t *testing.T) {
			fetchResponse, err := awxClient.CredentialsService.GetCredentialsByID(createResponse.ID, map[string]string{})
			if err != nil {
				t.Error(err)
			}

			if fetchResponse.Name != tt.data["name"] {
				t.Errorf("Expecting %s but got %s", tt.data["name"], fetchResponse.Name)
			}
		})

		t.Run("Update", func(t *testing.T) {
			tt.data["name"] = "credential_x"

			updateResponse, err := awxClient.CredentialsService.UpdateCredentialsByID(createResponse.ID, tt.data,
				map[string]string{})
			if err != nil {
				t.Error(err)
			}

			if updateResponse.Name != tt.data["name"] {
				t.Errorf("Expecting updated credentials name %s but got %s", tt.data["name"], updateResponse.Name)
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := awxClient.CredentialsService.DeleteCredentialsByID(createResponse.ID,
				map[string]string{})
			if err != nil {
				t.Error(err)
			}
		})
	}
}
