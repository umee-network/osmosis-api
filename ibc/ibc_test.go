package ibc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetIBCDenom(t *testing.T) {
	denom, err := GetIBCDenom("channel-0", "uatom")
	require.NoError(t, err)
	require.NotEmpty(t, denom)
	require.Equal(
		t,
		"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
		denom,
	)

	_, err = GetIBCDenom("", "")
	require.Error(t, err)
}
