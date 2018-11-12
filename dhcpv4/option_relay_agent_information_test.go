package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptRelayAgentInformation(t *testing.T) {
	data := []byte{
		byte(OptionRelayAgentInformation),
		13,
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}

	// short bytes
	opt, err := ParseOptRelayAgentInformation([]byte{})
	require.Error(t, err)

	// wrong code
	opt, err = ParseOptRelayAgentInformation([]byte{1, 2, 1, 0})
	require.Error(t, err)

	// wrong option length
	opt, err = ParseOptRelayAgentInformation([]byte{82, 3, 1, 0})
	require.Error(t, err)

	// short sub-option length
	opt, err = ParseOptRelayAgentInformation([]byte{82, 2, 1, 1})
	require.Error(t, err)

	opt, err = ParseOptRelayAgentInformation(data)
	require.NoError(t, err)
	require.Equal(t, len(opt.Options), 2)
	circuit, ok := opt.Options[0].(*OptionGeneric)
	require.True(t, ok)
	remote, ok := opt.Options[1].(*OptionGeneric)
	require.True(t, ok)
	require.Equal(t, circuit.Data, []byte("linux"))
	require.Equal(t, remote.Data, []byte("boot"))
}

func TestParseOptRelayAgentInformationToBytes(t *testing.T) {
	opt := OptRelayAgentInformation{}
	opt1 := &OptionGeneric{OptionCode: 1, Data: []byte("linux")}
	opt.Options = append(opt.Options, opt1)
	opt2 := &OptionGeneric{OptionCode: 2, Data: []byte("boot")}
	opt.Options = append(opt.Options, opt2)
	data := opt.ToBytes()
	expected := []byte{
		byte(OptionRelayAgentInformation),
		13,
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptRelayAgentInformationToBytesString(t *testing.T) {
	o := OptRelayAgentInformation{}
	require.Equal(t, "Relay Agent Information -> []", o.String())
}
