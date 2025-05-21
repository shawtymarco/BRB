import { Events, Client } from "discord.js";
import { deploy } from "../deploy";
import { Request } from "../api";

module.exports = {
    name: Events.ClientReady,
    once: true,
    async execute(client: Client) {
        await deploy();
        console.log(`Ready! Logged in as ${client.user?.tag}`);

        Request.get("hello");
    },
};