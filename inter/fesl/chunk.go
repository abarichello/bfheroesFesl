package fesl

// import (
// 	"github.com/Synaxis/bfheroesFesl/inter/network"
// 	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
// )

// const (
// 	Nextchunk = "GetNextChunk"
// )

// type ansChunk struct {
// 	TXN   string `fesl:"TXN"`
// 	Chunk string `fesl:"chunk"`
// 	next  string `fesl:"GetNextChunk"`
// }

// // func (fm *FeslManager) Chunk(event *network.EventNewClient) {
// // 	event.Client.Answer(&codec.Pkt{
// // 		Message: Chunk,
// // 		Content: ansChunk{
// // 			TXN:    chunk,
// // 			next:   "1",
// // 		},
// // 		Send: 0xC0000000,
// // 	})
// // }

// func (fm *FeslManager) Chunk(event network.EventClientProcess) {
// 	event.Client.Answer(&codec.Pkt{
// 		Content: ansChunk{TXN: Nextchunk},
// 		Message:    Chunk,
// 	})
// }
