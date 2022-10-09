package whois

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/zonedb/zonedb"
)

// Lookup performs a whois lookup to the appropriate whois-server using ZoneDB
func Lookup(domain string) ([]byte, error) {

	zone := zonedb.PublicZone(domain)
	if zone == nil {
		return nil, errors.New("could not identify TLDs corresponding whois-server")
	}

	whoisServer := fmt.Sprintf("%v:43", zone.WhoisServer())

	tcpAddr, err := net.ResolveTCPAddr("tcp4", whoisServer)
	if err != nil {
		return nil, fmt.Errorf("failed resolving tcp address: %w", err)
	}

	conn, err := net.Dial("tcp", tcpAddr.String())
	if err != nil {
		return nil, fmt.Errorf("failed connecting to whois server: %w", err)
	}
	defer conn.Close()

	payload := []byte(domain + "\r\n")
	_, err = conn.Write(payload)
	if err != nil {
		return nil, fmt.Errorf("failed requesting domain: %w", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed setting read deadline: %w", err)
	}

	reader := bufio.NewReader(conn)

	result, err := io.ReadAll(reader)
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			return nil, errors.New("request timed out")
		}
		if strings.Contains(err.Error(), "connection reset") {
			return nil, errors.New("request failed")
		}
		return nil, fmt.Errorf("failed reading response: %w", err)
	}

	return result, nil
}
