//This is lease support for nclient4

package nclient4

import (
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

//Lease contains a DHCPv4 lease after DORA.
//note: Lease doesn't include binding interface name
type Lease struct {
	ACK          *dhcpv4.DHCPv4
	CreationTime time.Time
	IDOptions    dhcpv4.Options //DHCPv4 options to identify the client like client-id, option82/remote-id
}

// WithClientIDOptions configures a list of DHCPv4 option code that DHCP server
// uses to identify client, beside the MAC address.
func WithClientIDOptions(cidl dhcpv4.OptionCodeList) ClientOpt {
	return func(c *Client) (err error) {
		c.clientIDOptions = cidl
		return
	}
}

//Release send DHCPv4 release messsage to server, based on specified lease.
//release is sent as unicast per RFC2131, section 4.4.4.
//This function requires assigned address has been added on the binding interface.
func (c *Client) Release(lease *Lease) error {
	if lease == nil {
		return fmt.Errorf("lease is nil")
	}
	req, err := dhcpv4.New()
	if err != nil {
		return err
	}
	//This is to make sure use same client identification options used during
	//DORA, so that DHCP server could identify the required lease
	req.Options = lease.IDOptions

	req.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeRelease))
	req.ClientHWAddr = lease.ACK.ClientHWAddr
	req.ClientIPAddr = lease.ACK.YourIPAddr
	req.UpdateOption(dhcpv4.OptServerIdentifier(lease.ACK.ServerIPAddr))
	req.SetUnicast()
	luaddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", lease.ACK.YourIPAddr, 68))
	if err != nil {
		return err
	}

	uniconn, err := net.DialUDP("udp4", luaddr, &net.UDPAddr{IP: lease.ACK.ServerIPAddr, Port: 67})
	if err != nil {
		return err
	}
	_, err = uniconn.Write(req.ToBytes())
	if err != nil {
		return err
	}
	c.logger.PrintMessage("sent message:", req)
	return nil
}
