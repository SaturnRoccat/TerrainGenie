import { world, system } from '@minecraft/server';
import { commandBuff, chunkPositions, eachChunkSize } from './levelData.js';


function dateTime() { return globalThis.__date_clock() }

let lastIndex = 0;
let startTime = 0
let lastChunk = 0

const realChunkSize = 256 // Hard Coded for now should be set by the level data

function placeBlockData(player) {
    startTime = dateTime();
    const tpPlayerToLatestChunkInTick = (player) => {
        if (lastIndex >= commandBuff.length) return false;
        if (lastIndex % (realChunkSize ) !== 0) return false;
        const chunkIndex = Math.floor(lastIndex / realChunkSize) ;
        if (lastChunk === chunkIndex) return false;
        lastChunk = chunkIndex;
        const chunkPosX = chunkPositions[chunkIndex * 2];
        const chunkPosZ = chunkPositions[chunkIndex * 2 + 1];
        system.run(() => player.teleport({ x: chunkPosX, y: 100, z: chunkPosZ }));
        return true;
    }
    for (let i = lastIndex; i < commandBuff.length; i++) {
        const command = commandBuff[i];
        system.run(() => player.runCommand(command));
        lastIndex = i;
        if (dateTime() - startTime > 180000 || tpPlayerToLatestChunkInTick(player)) {
            const callLambda = () => placeBlockData(player);
            tpPlayerToLatestChunkInTick(player);
            system.run(callLambda)
            break;
        }
    }

}

world.beforeEvents.chatSend.subscribe((eventData) => {
    const player = eventData.sender;
    switch (eventData.message) {
        case '!import_world_data':
            eventData.cancel = true;
            placeBlockData(player);
            break;
        default: break;
    }
});