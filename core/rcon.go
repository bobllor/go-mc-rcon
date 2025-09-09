package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"slices"
	"time"
)

type RCON struct {
	net.Conn
	Packet          RCONPacket
	IsAuthenticated bool
}

type RCONPacket struct {
	Length    int32
	RequestID int32
	Type      int32
	Payload   []byte
	Pad       []byte
}

// Connect intiiates a connection to the host.
func Connect(networkType string, host string, timeout ...time.Duration) (*RCON, error) {
	if len(timeout) > 1 || len(timeout) < 1 {
		errMsg := "timeout can only take one argument"
		return nil, errors.New(errMsg)
	}

	conn, err := net.DialTimeout(networkType, host, timeout[0])
	if err != nil {
		return nil, err
	}

	return &RCON{Conn: conn}, nil
}

// Authenticates authenticates the connection to the given RCON address. This is required to be called
// before sending a command to the server.
//
// It returns an error if the authentication process fails.
func (r *RCON) Authenticate(password string) error {
	if r.IsAuthenticated {
		fmt.Println("Already authenticated")
		return nil
	}

	packetType := int32(3)
	payload := []byte(password)

	r.Packet = RCONPacket{
		Length:    calculateLength(payload),
		RequestID: generateRequestID(),
		Type:      packetType,
		Payload:   payload,
		Pad:       []byte{0, 0},
	}

	buffer, err := r.buildBuffer()
	if err != nil {
		return err
	}

	_, err = r.transmitPacket(buffer)
	if err != nil {
		return err
	}

	r.IsAuthenticated = true
	return nil
}

// Command sends a command to the RCON server and the server's response as a string.
//
// The server's response can be empty, and will return an empty string.
// An error will be returned if any issues occur during the process.
func (r *RCON) Command(command string) (string, error) {
	if !r.IsAuthenticated {
		errMsg := "not authenticated, call RCON.Authenticate"
		return "", errors.New(errMsg)
	}
	packetType := int32(2)
	payload := []byte(command)

	r.Packet = RCONPacket{
		Length:    calculateLength(payload),
		RequestID: generateRequestID(),
		Type:      packetType,
		Payload:   payload,
		Pad:       []byte{0, 0},
	}

	buffer, err := r.buildBuffer()
	if err != nil {
		return "", err
	}

	response, err := r.transmitPacket(buffer)
	if err != nil {
		return "", err
	}

	payloadResponse := (*response)[12:]
	// ending index of the response up to the padding bytes.
	responseEndIndex := slices.IndexFunc(payloadResponse, func(i byte) bool { return i == 0 })

	// some commands will have no output, and the response will be 0s only.
	if responseEndIndex == 0 {
		return "", nil
	}
	responseStr := string(payloadResponse[:responseEndIndex-1])

	return responseStr, nil
}

// buildBuffer builds the buffer to be sent to the active connection.
func (r *RCON) buildBuffer() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	writeToEndianBuffer(buf, r.Packet.Length)
	writeToEndianBuffer(buf, r.Packet.RequestID)
	writeToEndianBuffer(buf, r.Packet.Type)

	_, err := buf.Write(r.Packet.Payload)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(r.Packet.Pad)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// transmitPacket sends the buffer to the active connection.
// It returns a pointer of the server's response in bytes and an error.
func (r *RCON) transmitPacket(buf *bytes.Buffer) (*[]byte, error) {
	_, err := r.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	response := make([]byte, 4096)
	_, err = r.Read(response)
	if err != nil {
		return nil, err
	}

	packetID := int32(binary.LittleEndian.Uint32(buf.Bytes()[4:8]))
	// -1 fails for any command, not just auth
	requestIDResponse := int32(binary.LittleEndian.Uint32(response[4:8]))
	if requestIDResponse == -1 {
		errMsg := fmt.Sprintf("command failure: %v", requestIDResponse)
		return nil, errors.New(errMsg)
	} else if requestIDResponse != packetID {
		// not sure if this error will occur but just in case?
		errMsg := fmt.Sprintf("failure to validate response ID: got %v, expected %v",
			requestIDResponse, packetID)
		return nil, errors.New(errMsg)
	}

	return &response, nil
}

// calculateLength calculates the total length of the packet.
// It requires the payload as an argument to calculate its final length.
func calculateLength(payload []byte) int32 {
	length := int32(4 + 4 + len(payload) + 2)

	return length
}

// generateRequestID generates a random int32 ID.
func generateRequestID() int32 {
	randomSrc := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := randomSrc.Int31()

	return id
}

// writeToEndianBuffer writes the given data to LittleEndian.
func writeToEndianBuffer(b *bytes.Buffer, data any) {
	_ = binary.Write(b, binary.LittleEndian, data)
}
