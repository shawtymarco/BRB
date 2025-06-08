import { CommandInteraction, SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
    .setName("cmd")
    .setDescription("to execute in-game commands");

export async function execute(interaction: CommandInteraction) {
    return interaction.reply("WIP!");
}
