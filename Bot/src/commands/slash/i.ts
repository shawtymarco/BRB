import { CommandInteraction, SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
    .setName("i")
    .setDescription("To view player stats");

export async function execute(interaction: CommandInteraction) {
    return interaction.reply("WIP");
}
