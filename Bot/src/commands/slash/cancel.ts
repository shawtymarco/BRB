import { ActionRowBuilder, ButtonBuilder, ButtonStyle, ChannelType, CommandInteraction, MessageFlags, PrivateThreadChannel, SlashCommandBuilder } from "discord.js";
import { EmbedUtil } from "../../core/EmbedUtil";
import { gamesDB } from "../../core/GameCore";

export const data = new SlashCommandBuilder()
    .setName("cancel")
    .setDescription("Start a void vote (game thread only, once)");
// in complete, no test yet
export async function execute(interaction: CommandInteraction) {
    if (!interaction.channel?.isThread() || interaction.channel.type !== ChannelType.PrivateThread) {
        await interaction.reply({
            embeds: [EmbedUtil.create({
                type: 'no',
                description: 'You can only use this command in a game\'s private thread.'
            })],
            flags: MessageFlags.Ephemeral
        });
        return;
    }

    const thread = interaction.channel as PrivateThreadChannel;
    const game = gamesDB.data.get(thread.id);
    if (!game) {
        await interaction.reply({
            embeds: [EmbedUtil.create({ type: 'no', description: 'No active game found for this thread.' })],
            flags: MessageFlags.Ephemeral
        });
        return;
    }

    if (game.cancelVoteStarted) {
        await interaction.reply({
            embeds: [EmbedUtil.create({ type: 'no', description: 'Cancel vote already started for this game.' })],
            flags: MessageFlags.Ephemeral
        });
        return;
    }

    game.cancelVoteStarted = true;
    game.cancelVoteAgreeUserIds = [];
    game.cancelVoteDisagreeUserIds = [];

    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(
        new ButtonBuilder().setCustomId('cancel_yes').setLabel('Yes').setStyle(ButtonStyle.Success),
        new ButtonBuilder().setCustomId('cancel_no').setLabel('No').setStyle(ButtonStyle.Danger),
    );

    const msg = await thread.send({
        content: (game as any)["memberIds"].map((id: string) => `<@${id}>`).join(" "),
        embeds: [EmbedUtil.create({
            type: 'yes',
            description: 'Vote to void the current game. Required agreements: 3 (2v2), 4 (3v3), 6 (4v4).'
        })],
        components: [row]
    });

    game.cancelVoteMessageId = msg.id;
    await gamesDB.save();

    await interaction.reply({
        embeds: [EmbedUtil.create({ type: 'yes', description: 'Cancel vote started.' })],
        flags: MessageFlags.Ephemeral
    });
}


