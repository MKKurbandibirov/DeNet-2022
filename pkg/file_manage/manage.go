package filemanage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"wallet/pkg/encrypt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrPermDenied = "permission denied"

func GetPassword() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	password := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return password, nil
}

func GenPassword() (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)	

	password := hexutil.Encode(privateKeyBytes)[2:]

	fmt.Printf("Your new password: %s\n", password)
	return password, nil
}

func NewClientHandling(file *os.File) error {
	fmt.Println("You are a new client, so I generate your new password!")
	password, err := GenPassword()
	if err != nil {
		return err
	}

	hash, err := encrypt.GenHash(password)
	if err != nil {
		return err
	}

	cypherText, err := encrypt.Encrypt(hash)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprint(file, cypherText, "\n"); err != nil {
		return err
	}
	return nil
}

func ClientHandling(file *os.File) error {
	scanner := bufio.NewScanner(file)

	fmt.Print("Enter your password: ")
	password, err := GetPassword()
	if err != nil {
		return err
	}

	hash, err := encrypt.GenHash(password)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		text, err := encrypt.Decrypt(scanner.Text())
		if err != nil {
			return err
		}

		if text[:16] == hash[:16] {
			return nil
		}
	}

	return errors.New(ErrPermDenied)
}

func Access(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	defer file.Close()

	if err != nil {
		if os.IsNotExist(err) {
			newfile, err := os.Create(fileName)
			if err != nil {
				return err
			}
			return NewClientHandling(newfile)
		}
		return err
	}

	if err = ClientHandling(file); err != nil && err.Error() == ErrPermDenied {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Are you a new client(yes/no): ")
		for scanner.Scan() {
			switch scanner.Text() {
			case "yes":
				return NewClientHandling(file)
			case "no":
				return err
			}
			fmt.Println("\nType yes/no only!")
			fmt.Print("Are you a new client(yes/no): ")
		}
	}

	return nil
}
