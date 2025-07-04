import { ChannelType, Events, GuildMember, Interaction, MessageFlags } from "discord.js";
import { gamesDB } from "../core/GameCore";
import { EmbedUtil } from "../core/EmbedUtil";

module.exports = {
    name: Events.InteractionCreate,
    async execute(interaction: Interaction) {
        if (!interaction.isStringSelectMenu()) return;

        if (interaction.channel?.isThread() && interaction.channel.type === ChannelType.PrivateThread) {
            const thread = interaction.channel;
            if (interaction.customId === "pick_teammates") {
                const game = gamesDB.data.get(thread.id);
                if (game != null) {
                    const captains = await game.captains();
                    if (game.isTeam1Turn() && interaction.user.id === captains[0].id || !game.isTeam1Turn() && interaction.user.id === captains[1].id) {
                        interaction.values.forEach(id => {
                            (game.isTeam1Turn() ? game.team1Ids : game.team2Ids).push(id)
                        });
                        game.step++;
                        await game.updateCaptainPickingMessage();

                        interaction.reply({
                            embeds: [EmbedUtil.create({
                                type: "yes",
                                description: `You have picked ${interaction.values.map(id => `<@${id}>`).join(' & ')}. It is the other captain's turn to pick now.`
                            })], flags: MessageFlags.Ephemeral
                        });
                    } else if (interaction.user.id === captains[0].id || interaction.user.id === captains[1].id) {
                        await game.updateCaptainPickingMessage();
                        interaction.reply({embeds: [EmbedUtil.create({
                            type: "no",
                            description: "Please wait your turn to pick a teammate"
                        })], flags: MessageFlags.Ephemeral});
                    } else {
                        await game.updateCaptainPickingMessage();
                        interaction.reply({
                            embeds: [EmbedUtil.create({
                                type: "no",
                                description: "You are not allowed to pick teammates. Please wait for your captain to pick!"
                            })], flags: MessageFlags.Ephemeral
                        });
                    }
                }
            }
        }
    },
};
