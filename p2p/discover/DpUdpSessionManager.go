package discover

import (
	"errors"
	"github.com/xtaci/kcp-go"
	"net"
	"sync"
)

const dataShards = 10
const parityShards = 3

type UdpSessionManager interface {
	Accept() (*DpUdpSession, error)
	Close() error
	LocalAddr() net.Addr
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
}

type DpUdpSessionManager struct {
	mutex sync.Mutex

	baseListener *kcp.Listener

	sessions map[string]*DpUdpSession
}

func (sm *DpUdpSessionManager) Close() error {
	if sm.baseListener != nil {
		return sm.baseListener.Close()
	}

	return nil
}

func (sm *DpUdpSessionManager) Accept() (*DpUdpSession, error) {
	if sm.baseListener == nil {
		return nil, errors.New("listener not initialized")
	}

	conn, err := sm.baseListener.Accept()

	if err != nil {
		return nil, err
	}

	udpAddr := *conn.RemoteAddr().(*net.UDPAddr)
	remoteAddr := udpAddr.IP.String() + ":30303" //todo
	remoteUdpAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		return nil, err
	}

	dpUdpSession := &DpUdpSession{
		BaseConn: conn,
		addr:     remoteUdpAddr,
	}

	return dpUdpSession, nil
}

func (sm *DpUdpSessionManager) LocalAddr() net.Addr {
	return sm.baseListener.Addr()
}

func (sm *DpUdpSessionManager) WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error) {
	toAddr := addr.IP.String() + ":30303" //todo
	if session, err := createDpUdpSession(toAddr); err == nil {
		//	sm.sessions[toAddr] = session
		return session.Write(b)
	} else {
		return 0, err
	}
}

func CreateDpUdpSessionManager(address string) (*DpUdpSessionManager, error) {
	if len(address) == 0 {
		sessionManager := &DpUdpSessionManager{
			baseListener: nil,
			sessions:     make(map[string]*DpUdpSession),
		}
		return sessionManager, nil
	}

	if listener, err := kcp.ListenWithOptions(address, nil, dataShards, parityShards); err == nil {
		sessionManager := &DpUdpSessionManager{
			baseListener: listener,
			sessions:     make(map[string]*DpUdpSession),
		}
		return sessionManager, nil
	} else {
		return nil, err
	}
}

func createDpUdpSession(targetAddress string) (s *DpUdpSession, err error) {
	if session, err := kcp.DialWithOptions(targetAddress, nil, dataShards, parityShards); err == nil {
		session.SetMtu(1280)
		dpUdpSession := &DpUdpSession{
			BaseConn: session,
			addr:     session.RemoteAddr().(*net.UDPAddr),
		}

		return dpUdpSession, nil
	} else {
		return nil, err
	}
}
