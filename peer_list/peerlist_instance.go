package peer_list

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/types"
)

const (
	DataReq            = "DataReq"     //DATA STORE SEND TO DATA PROVIDER TO GET DATA
	DataRecv           = "DataRecv"    //MSG SEND TO PROVIDER BY PEERLIST TO INFORM THE END OF A FILE TRANSPORT
	MetaDataReq        = "MetaDataReq" //SLAVE SEND IT TO MASTER TO GET METADATA
	MetaDataRes        = "MetaDataRes"
	InfoReq            = "TrackReq" //SEND TO TRACKER TO GET DATA POSITION
	InfoRes            = "InfoRes"
	SeekReceiverReq    = "SeekReceiverReq" //DATA PROVIDER SEND TO SEEK FOR RECEIVER
	SeekParticipantReq = "SeekParticipantReq"
)

type PeerListInstance struct {
	ID       string
	PL       *PeerList
	channels map[string]chan MessageInfo
	bus      chan []byte
}

func NewPeerListInstance(id string, pl *PeerList) *PeerListInstance {
	return &PeerListInstance{
		ID:       id,
		PL:       pl,
		channels: make(map[string]chan MessageInfo),
		bus:      make(chan []byte, 10),
	}
}

func (this *PeerListInstance) HandleMsg(meth string, msg MessageInfo) {
	method := this.channels[meth]
	if method != nil {
		method <- msg
	}
}

func (this *PeerListInstance) AddChannel(name string, ch chan MessageInfo) {
	this.channels[name] = ch
}
func (this *PeerListInstance) RemoteProcedureCall(to types.Address, method string, msg []byte) error {
	pkg := newMessage(this.PL.Address, to, this.ID, method, msg, 1)
	pkgb, err := json.Marshal(pkg)
	if err != nil {
		return err
	}
	this.PL.netService.SendMessageTo(to.ToString(), pkgb)
	return nil
}

func (this *PeerListInstance) SendDataPackTo(to types.Address, pack types.DataPack) {

}

func (this *PeerListInstance) UpdatePunishedPeers(map[string][]types.Address) {

}

func (this *PeerListInstance) Close() {
	this.PL.RemoveInstance(this.ID)
	for k, v := range this.channels {
		close(v)
		delete(this.channels, k)
	}
}
