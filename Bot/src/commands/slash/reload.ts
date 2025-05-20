import { CommandInteraction } from "discord.js";

import { SlashCommandBuilder } from 'discord.js';
import { slashCommands } from "../..";

export const data = new SlashCommandBuilder()
    .setName('reload')
    .setDescription('Reloads a command.')
    .addStringOption(option =>
        option.setName('command')
            .setDescription('The command to reload.')
            .setRequired(true));

export async function execute(interaction: CommandInteraction) {
    const commandName = interaction.options.get('command', true).name.toLowerCase();
    const command = slashCommands.get(commandName);

    if (!command) {
        return interaction.reply(`There is no command with name \`${commandName}\`!`);
    }

    delete require.cache[require.resolve(`./${command.data.name}.js`)];

    try {
        const newCommand = require(`./${command.data.name}.js`);
        slashCommands.set(newCommand.data.name, newCommand);
        await interaction.reply(`Command \`${newCommand.data.name}\` was reloaded!`);
    } catch (error: any) {
        console.error(error);
        await interaction.reply(`There was an error while reloading a command \`${command.data.name}\`:\n\`${error.message}\``);
    }
}