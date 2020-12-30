package whois

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/zonedb/zonedb"
)

// Lookup performs a whois lookup to the appropriate whois-server using ZoneDB
func Lookup(domain string) ([]byte, error) {

	zone := zonedb.PublicZone(domain)
	if zone == nil {
		return nil, errors.New("Could not identify TLDs corresponding whois-server")
	}

	whoisServer := fmt.Sprintf("%v:43", zone.WhoisServer())

	tcpAddr, err := net.ResolveTCPAddr("tcp4", whoisServer)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", tcpAddr.String())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload := []byte(domain + "\r\n")
	_, err = conn.Write(payload)
	if err != nil {
		return nil, err
	}

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return result, nil
}
