package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	ibcPrefix      = "ibc"
	transferPrefix = "transfer"
)

func GetIBCDenom(channelID string, denom string) (string, error) {
	if len(channelID) == 0 || len(denom) == 0 {
		return "", fmt.Errorf("channel ID and denom cannot be empty")
	}

	hash := tmhash.New()
	sourceStr := fmt.Sprintf("%s/%s/%s", transferPrefix, channelID, denom)

	if _, err := hash.Write([]byte(sourceStr)); err != nil {
		return "", err
	}

	bz := hash.Sum(nil)
	hashString := strings.ToUpper(hex.EncodeToString(bz))

	return fmt.Sprintf("%s/%s", ibcPrefix, hashString), nil
}
