package main

import (
	"fmt"

	"github.com/MatthewAraujo/notify/encrypt"
)

func main() {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY3NzA0MTQsImlhdCI6MTcxNjc2OTgxNCwiaXNzIjoiODk3MzkyIn0.cOL2Jp2x0BHkQLY9pfPKZtg6pR9fwOdEKNg-VLzya6NsVHcTfQbmpn3Xx7T0ZF9e04Gv3WDSDpYFVxrw_aKeA_NzD-9P-nf6nuJehh6eCODzzsJ14AzoHfIbVP6_qo8e9zCuKQEGo47PM6ZAvtqZDhpyVIvGNGZjkyblKrNwkYxYk2Xuy5Rsgp2xrRaWDxy0_kZhebI18UszK1J3WPUVEpHdf5J-hv48IyvqsaLxMGG5TP-MK-eUcRD6DIhTUvjj5ZMLf3l6vB-d9OSFpVQDV6iKQUbWz8LGt8Q6WtZOdTekKV6nmQH913yJnM6I49IJZC9voA96dLUgNay3ethZSg"
	encryptedToken, err := encrypt.EncryptToken(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Encrypted token:", encryptedToken)
}
