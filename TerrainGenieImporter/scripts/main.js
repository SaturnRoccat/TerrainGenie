import { world, system } from "@minecraft/server";
import { commandBuff, chunkPositions, eachChunkSize } from "./levelData.js";

const dateTime = () => { return globalThis.__date_clock() };

let lastIndex = 0;
let startTime = 0;
let lastChunk = 0;

const realChunkSize = 256; // Hard Coded for now should be set by the level data

/** @param { import("@minecraft/server").Player } player */
const lastChunkInTick = () => {
    const chunkIndex = Math.floor(lastIndex / realChunkSize);
    if (
        lastIndex >= commandBuff.length
        || lastIndex % (realChunkSize) !== 0
        || lastChunk === chunkIndex
    ) return false;

    lastChunk = chunkIndex;
    return true;
};

/** @param { import("@minecraft/server").Player } player */
const placeBlockData = (player) => {
    startTime = dateTime();
    for (let i = lastIndex; i < commandBuff.length; i++) {
        const command = commandBuff[i];
        system.run(() => player.runCommand(command));

        lastIndex = i;
        if (
            (dateTime() - startTime) > 180000
            || lastChunkInTick()
        ) {
            const chunkPosX = chunkPositions[lastChunk * 2];
            const chunkPosZ = chunkPositions[lastChunk * 2 + 1];
            system.run(() => {
                player.teleport({ x: chunkPosX, y: 100, z: chunkPosZ });
                system.run(() => placeBlockData(player));
            });

            break;
        };
    };
};

world.beforeEvents.chatSend.subscribe((eventData) => {
    const { sender: player, message } = eventData;
    if (!message.startsWith("!")) return;
    eventData.cancel = true;

    const command = message.split(" ")[0].slice(1).trim();
    switch(command) {
        case "import_world_data": placeBlockData(player); break;
    };
});