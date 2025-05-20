import { Events, Client } from "discord.js";
import { deploy } from "../deploy";

module.exports = {
    name: Events.ClientReady,
    once: true,
    async execute(client: Client) {
        await deploy();
        console.log(`Ready! Logged in as ${client.user?.tag}`);
    },
};