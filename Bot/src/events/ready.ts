import { Events, Client } from "discord.js";
import { deploy } from "../deploy";
import { APIEndpoints, Request } from "../api";
import { Game, gamesDB } from "../core/GameCore";
import { dconfig } from "../config";

module.exports = {
    name: Events.ClientReady,
    once: true,
    async execute(client: Client) {
        await deploy();
        console.log(`Ready! Logged in as ${client.user?.tag}`);

        detectServerRestart();

        await gamesDB.load(raw => Object.assign(new Game(raw.id, raw.threadId, raw.vcId, "", "", raw.teamSize, raw.memberIds, raw.captainIds), raw));
    },
};

function detectServerRestart() {
    let initialTime: number | null = null;
    setInterval(async () => {
        try {
            const res = await Request.get(APIEndpoints.CONNECT);
            const currentTime = res.time;

            if (initialTime === null) {
                initialTime = currentTime;
            } else if (currentTime !== initialTime) {
                console.log("⚠️ Server restart detected!");
                gamesDB.data.map(game => {
                    game.terminateGame();
                })
                initialTime = currentTime;
            }
        } catch (err) {
            console.error(`Failed to connect to ${dconfig.api}`);
        }
    }, 1000);
}