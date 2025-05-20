import { Events, Message } from "discord.js";
import { messageCommands } from "..";
import { config } from "../config";

export default {
    name: Events.MessageCreate,
    async execute(message: Message) {
        if (message.author.bot || !message.content.startsWith(config.prefix)) return;
        const args = message.content.slice(config.prefix.length).trim().split(/ +/);
        const commandName = args.shift()?.toLowerCase();
        if (!commandName) return;
        const command = messageCommands.get(commandName);
        if (!command) return;
        try {
            await command.execute(message, args);
        } catch (err) {
            console.error(err);
            message.reply("Error running command.");
        }
    },
};
