import { REST, Routes } from "discord.js";
import { slashCommands } from ".";
import { token, config } from "./config";

const rest = new REST().setToken(token);

export async function deploy() {
	try {
		console.log(`Started refreshing ${slashCommands.size} application (/) commands.`);

		const data: any = await rest.put(
            Routes.applicationGuildCommands(config.clientId, config.guildId),
			{ body: [...slashCommands.values()].map(cmd => cmd.data.toJSON()) }
		);

		console.log(`Successfully reloaded ${data.length} application (/) commands.`);
	} catch (error) {
		console.error(error);
	}
}