import { Events, Client } from "discord.js";
import { deploy } from "../deploy";
import { APIEndpoints, Request } from "../api";

module.exports = {
    name: Events.ClientReady,
    once: true,
    async execute(client: Client) {
        await deploy();
        console.log(`Ready! Logged in as ${client.user?.tag}`);

        const data = await Request.get(APIEndpoints.CONNECT);
        console.log(data.message);
    },
};