import {Client, Collection, GatewayIntentBits, MessageFlags} from "discord.js";
import {token} from "./config";
import path from "path";
import fs from "fs";
import {pathToFileURL} from "url";

export let slashCommands = new Collection<string, any>();
export let messageCommands = new Collection<string, any>();

export var client = new Client({
    intents: [
        GatewayIntentBits.Guilds,
        GatewayIntentBits.GuildMessages,
        GatewayIntentBits.DirectMessages,
        GatewayIntentBits.MessageContent,
        GatewayIntentBits.GuildVoiceStates
    ],
});

(async () => {
    // Load SLASH commands
    const slashPath = 'src/commands/slash';
    const slashFiles = fs.readdirSync(slashPath).filter(file => file.endsWith('.ts'));

    for (const file of slashFiles) {
        const filePath = path.join(slashPath, file);
        const command = await import(pathToFileURL(filePath).href);
        if ('data' in command && 'execute' in command) {
            slashCommands.set(command.data.name, command);
        }
    }

    // Load MESSAGE commands
    const msgPath = 'src/commands/message';
    const msgFiles = fs.readdirSync(msgPath).filter(file => file.endsWith('.ts'));

    for (const file of msgFiles) {
        const filePath = path.join(msgPath, file);
        const command = await import(pathToFileURL(filePath).href);
        if ('name' in command && 'execute' in command) {
            messageCommands.set(command.name, command);
        }
    }

    // Load EVENTS
    const eventsPath = 'src/events';
    const eventFiles = fs.readdirSync(eventsPath).filter(file => file.endsWith('.ts'));

    for (const file of eventFiles) {
        const filePath = path.join(eventsPath, file);
        const event = (await import(pathToFileURL(filePath).href)).default;
        if (event.once) {
            client.once(event.name, (...args) => event.execute(...args));
        } else {
            client.on(event.name, (...args) => event.execute(...args));
        }
    }
})();

client.login(token);

// InitiateWebsocket();