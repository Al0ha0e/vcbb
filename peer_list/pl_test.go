package peer_list

import (
	"fmt"
	"testing"
	"vcbb/net"

	"github.com/Al0ha0e/vcbb/types"
)

func TestPeerList(t *testing.T) {
	ns := net.NewNetSimulator()
	var addrs [3]types.Address
	addrs[0] = types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addrs[1] = types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
	addrs[2] = types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d3")
	var pls [3]*PeerList
	for i := 0; i < 3; i++ {
		pls[i] = NewPeerList(addrs[i], ns)
	}
	pls[0].peers = append(pls[0].peers, []types.Address{addrs[1], addrs[2]}...)
	pls[1].peers = append(pls[1].peers, []types.Address{addrs[0], addrs[2]}...)
	pls[2].peers = append(pls[2].peers, []types.Address{addrs[0], addrs[1]}...)
	var chs [3]chan MessageInfo
	for i := 0; i < 2; i++ {
		chs[i] = make(chan MessageInfo, 1)
		pls[i].AddChannel("test", chs[i])
		go func(id int) {
			for {
				msg := <-chs[id]
				fmt.Println(id, msg.From.ToString(), msg.Content)
			}
		}(i)
		pls[i].Run()
	}
	pls[0].RemoteProcedureCall(addrs[1], "test", []byte{1, 2, 3})
	pls[1].RemoteProcedureCall(addrs[0], "test", []byte{4, 5, 6})
	//pls[2].BroadCastRPC("test", []byte{9, 9, 9, 0}, 3)\
	var sess [2]*PeerListInstance
	var chh [2]chan MessageInfo
	for i := 0; i < 2; i++ {
		chh[i] = make(chan MessageInfo, 1)
		sess[i] = pls[i].GetInstance("mmm")
		sess[i].AddChannel("test2", chh[i])
		go func(id int) {
			for {
				msg := <-chh[id]
				fmt.Println(id, msg.From.ToString(), msg.Content)
			}
		}(i)
	}
	sess[0].RemoteProcedureCall(addrs[1], "test2", []byte{1, 1, 1})
	sess[1].RemoteProcedureCall(addrs[0], "test2", []byte{3, 3, 3})
	for {
	}
}
