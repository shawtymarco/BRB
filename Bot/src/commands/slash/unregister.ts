import { CommandInteraction, SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
    .setName("unregister")
    .setDescription("To unregister your Discord account from your MC account");

export async function execute(interaction: CommandInteraction) {
    return interaction.reply("WIP");
}
