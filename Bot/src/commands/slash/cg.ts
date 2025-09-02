import { ChannelType, CommandInteraction, MessageFlags, PrivateThreadChannel, SlashCommandBuilder } from "discord.js";
import { EmbedUtil } from "../../core/EmbedUtil";
import { gamesDB } from "../../core/GameCore";

export const data = new SlashCommandBuilder()
    .setName("cg")
    .setDescription("to close a game's VC and thread");

export async function execute(interaction: CommandInteraction) {
    if (!interaction.channel?.isThread() || interaction.channel.type !== ChannelType.PrivateThread) {
        interaction.reply({
            embeds: [EmbedUtil.create({
                type: 'no',
                description: 'You can only execute this command in a game\'s private thread.'
            })]
        });
        return;
    }
    await interaction.deferReply({ flags: MessageFlags.Ephemeral });

    const thread = interaction.channel as PrivateThreadChannel;
    const game = gamesDB.data.get(thread.id);
    if (game != null) {
        await game.terminateGame(game);
    }
}
