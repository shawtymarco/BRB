import { Events, Client } from "discord.js";
import { deploy } from "../deploy";
import { APIEndpoints, Request } from "../api";
import { Game, gamesDB } from "../core/GameCore";

module.exports = {
    name: Events.ClientReady,
    once: true,
    async execute(client: Client) {
        await deploy();
        console.log(`Ready! Logged in as ${client.user?.tag}`);

        const data = await Request.get(APIEndpoints.CONNECT);
        console.log(data.message);

        await gamesDB.load(raw => Object.assign(new Game(raw.id, raw.threadId, raw.vcId, raw.teamSize, raw.memberIds, raw.captainIds), raw));
    },
};