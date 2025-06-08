import { CommandInteraction, SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
    .setName("lb")
    .setDescription("To view several types of leaderboards");

export async function execute(interaction: CommandInteraction) {
    return interaction.reply("WIP");
}
