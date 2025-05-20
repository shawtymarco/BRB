import { CommandInteraction, SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
    .setName("cmd")
    .setDescription("to uexecute in-game commands");

export async function execute(interaction: CommandInteraction) {
    return interaction.reply("WIP!");
}
