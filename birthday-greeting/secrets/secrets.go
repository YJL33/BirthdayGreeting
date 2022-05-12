package secrets

import (
	// "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Secret struct {
	UserName            string `json:"username"`
	Password            string `json:"password"`
	Engine              string `json:"engine"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	SSL                 bool   `json:"ssl"`
	DBName              string `json:"dbname"`
	DBClusterIdentifier string `json:"dbClusterIdentifier"`
}

func GetSecret() (Secret, error) {
	secretName := "rdsSecret"
	region := "ap-northeast-1"

	//Create a Secrets Manager client
	sess, err := session.NewSession()
	if err != nil {
		// Handle session creation error
		fmt.Println(err.Error())
		return Secret{}, err
	}
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// We only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			fmt.Println(err.Error())
		}
		return Secret{}, err
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	// var secretString, decodedBinarySecret string
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	}
	// } else {
	// 	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	// 	len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	// 	if err != nil {
	// 		fmt.Println("Base64 Decode Error:", err)
	// 		return Secret{}, err
	// 	}
	// 	decodedBinarySecret = string(decodedBinarySecretBytes[:len])
	// }

	secret := Secret{}
	json.Unmarshal([]byte(secretString), &secret)
	return secret, nil
}
