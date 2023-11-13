import { world, system } from '@minecraft/server';
import { commandBuff, RLE } from './levelData.js';


function dateTime() {return globalThis.__date_clock()}

let lastIndex = 0;
let startTime = 0


function placeBlockData(player) {
    startTime = dateTime();

    for (let i = lastIndex; i < commandBuff.length; i++) {
        const command = commandBuff[i];
        system.run(() => player.runCommand(command));
        lastIndex = i;
        if (dateTime() - startTime > 200000) {
            console.log(`Last index: ${lastIndex}`);
            const callLambda = () => placeBlockData(player);
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