import { Events, Message } from "discord.js";
import { messageCommands } from "..";
import { dconfig } from "../config";

export default {
    name: Events.MessageCreate,
    async execute(message: Message) {
        const isCommand = await checkCommandExecution(message);

        if (!isCommand && message.channel.id === dconfig.channels.register && !message.author.bot) message.delete().catch(() => {});
    },
};

async function checkCommandExecution(message: Message) {
    if (message.author.bot || !message.content.startsWith(dconfig.prefix)) return false;
    const args = message.content.slice(dconfig.prefix.length).trim().split(/ +/);
    const commandName = args.shift()?.toLowerCase();
    if (!commandName) return false;
    const command = messageCommands.get(commandName);
    if (!command) return false;
    try {
        await command.execute(message, args);
    } catch (err) {
        console.error(err);
        message.reply("Error running command.");
    }
    return true;
} 