package main

import (
	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/protocol/event"
)

func translateToCommands(commandBuff *[]string, PalletData *TG_Pallet_Data) {
	for _, chunk := range WorldChunks {
		for x := int32(0); x < ChunkWidth; x++ {
			for z := int32(0); z < ChunkDepth; z++ {
				var arrayBuf []byte
				createFillBlockCommand(&arrayBuf, uint32(x+chunk.ChunkPosition.x*ChunkWidth), uint32(chunk.HeightMap[x+z*ChunkWidth]), uint32(z+chunk.ChunkPosition.z*ChunkDepth), chunk.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, int32(chunk.HeightMap[x+z*ChunkWidth]), z})], PalletData)

				*commandBuff = append(*commandBuff, string(arrayBuf))
			}
		}
	}
}

func enableWSServer(PalletData *TG_Pallet_Data) {
	var server = mcwss.NewServer(nil)

	server.OnConnection(func(player *mcwss.Player) {
		println("Player connected")
		player.OnPlayerMessage(func(message *event.PlayerMessage) {
			if message.Message == "!set_world_data" {
				var commandBuff []string
				translateToCommands(&commandBuff, PalletData)
				for _, command := range commandBuff {
					player.ExecAs(command, func(statusCode int) {})
				}
			}
		})
	})

	println(server.Run().Error())

}
