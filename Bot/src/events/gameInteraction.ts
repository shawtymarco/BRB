import { ActionRowBuilder, ButtonInteraction, ChannelType, Events, Interaction, MessageFlags, StringSelectMenuBuilder, StringSelectMenuInteraction, StringSelectMenuOptionBuilder } from "discord.js";
import { gamesDB } from "../core/GameCore";
import { EmbedUtil } from "../core/EmbedUtil";

module.exports = {
    name: Events.InteractionCreate,
    async execute(interaction: Interaction) {
        if (interaction.channel?.isThread() && interaction.channel.type === ChannelType.PrivateThread) {
            const thread = interaction.channel;
            if (interaction.isStringSelectMenu() && interaction.customId === "pick_teammates") {
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
                    } else if (interaction.values.some(id => game.team1Ids.includes(id) || game.team2Ids.includes(id))) {
                        await game.updateCaptainPickingMessage();
                        interaction.reply({
                            embeds: [EmbedUtil.create({
                                type: "no",
                                description: "This player has already been picked. Pick someone else!"
                            })], flags: MessageFlags.Ephemeral
                        });
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
            } else if (interaction.isButton() && (interaction.customId === 'mapvote_yes' || interaction.customId === 'mapvote_no' || interaction.customId === 'cancel_yes' || interaction.customId === 'cancel_no')) {
                const game = gamesDB.data.get(thread.id);
                if (!game) return;
                const isMapVote = interaction.customId.startsWith('mapvote');
                const yes = interaction.customId.endsWith('yes');

                const agreeList = isMapVote ? game.mapVoteAgreeUserIds : game.cancelVoteAgreeUserIds;
                const disagreeList = isMapVote ? game.mapVoteDisagreeUserIds : game.cancelVoteDisagreeUserIds;

                const removeFrom = yes ? disagreeList : agreeList;
                const addTo = yes ? agreeList : disagreeList;
                const idx = removeFrom.indexOf(interaction.user.id);
                if (idx >= 0) removeFrom.splice(idx, 1);
                if (!addTo.includes(interaction.user.id)) addTo.push(interaction.user.id);

                const teamSize = game.teamSize;
                const thresholds: Record<number, number> = { 2: 3, 3: 4, 4: 6 };
                const required = thresholds[teamSize] ?? 3;

                const agreeCount = agreeList.length;
                const totalPlayers = (game as any).memberIds?.length ?? 0;

                const baseDesc = isMapVote ? `Map vote in progress. Required agreements: ${required}.` : `Void vote in progress. Required agreements: ${required}.`;
                const status = `Yes: ${agreeList.map((id: string) => `<@${id}>`).join(' ') || '—'} | No: ${disagreeList.map((id: string) => `<@${id}>`).join(' ') || '—'}`;

                const msgId = isMapVote ? game.mapVoteMessageId : game.cancelVoteMessageId;
                if (msgId) {
                    const msg = await thread.messages.fetch(msgId).catch(() => null);
                    if (msg) await msg.edit({ embeds: [EmbedUtil.create({ type: 'yes', description: `${baseDesc}\n${status}` })] });
                }

                await gamesDB.save();

                if (agreeCount >= required) {
                    if (isMapVote) {
                        const options = [
                            new StringSelectMenuOptionBuilder().setLabel('BW-Aquarium').setValue('BW-Aquarium'),
                            new StringSelectMenuOptionBuilder().setLabel('BW-Archway').setValue('BW-Archway'),
                            new StringSelectMenuOptionBuilder().setLabel('BW-Boletum').setValue('BW-Boletum'),
                            new StringSelectMenuOptionBuilder().setLabel('BW-Invasion').setValue('BW-Invasion'),
                            new StringSelectMenuOptionBuilder().setLabel('BW-Katsu').setValue('BW-Katsu'),
                            new StringSelectMenuOptionBuilder().setLabel('BW-Lectus').setValue('BW-Lectus'),
                            new StringSelectMenuOptionBuilder().setLabel('BW-Planet98').setValue('BW-Planet98'),
                        ];
                        const row = new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(
                            new StringSelectMenuBuilder().setCustomId('mapvote_select').setPlaceholder('Select a map').addOptions(options)
                        );
                        const selectMsg = await thread.send({
                            embeds: [EmbedUtil.create({ type: 'yes', description: 'Map vote accepted. Choose a map.' })],
                            components: [row]
                        });
                        game.mapSelectMessageId = selectMsg.id;
                        await gamesDB.save();
                    } 
                } else if (agreeCount === 0 && (agreeList.length + disagreeList.length) >= totalPlayers) {
                    await thread.send({ embeds: [EmbedUtil.create({ type: 'no', description: isMapVote ? 'Map vote rejected.' : 'Void vote rejected.' })] });
                }

                await (interaction as ButtonInteraction).reply({ embeds: [EmbedUtil.create({ type: 'yes', description: 'Your vote was recorded.' })], flags: MessageFlags.Ephemeral });

            } else if (interaction.isStringSelectMenu() && interaction.customId === 'mapvote_select') {
                const game = gamesDB.data.get(thread.id);
                if (!game) return;
                const chosen = interaction.values[0];
                (game as any).selectedMap = chosen;
                await gamesDB.save();
                await thread.send({ embeds: [EmbedUtil.create({ type: 'yes', description: `Selected map: ${chosen}. This map will be used for game creation.` })] });
                await (interaction as StringSelectMenuInteraction).reply({ embeds: [EmbedUtil.create({ type: 'yes', description: 'Map selection recorded.' })], flags: MessageFlags.Ephemeral });
            }
        }
    },
};
