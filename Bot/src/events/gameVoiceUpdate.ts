import {
    ChannelType,
    Collection,
    Events,
    Guild,
    GuildMember,
    OverwriteResolvable,
    OverwriteType,
    PermissionFlagsBits,
    TextChannel,
    VoiceChannel
} from "discord.js";
import { dconfig } from "../config";
import { APIEndpoints, Request } from "../api";
import { EmbedUtil } from "../core/EmbedUtil";
import { CacheUtil } from "../core/CacheUtil";
import { Game, gamesDB } from "../core/GameCore";

const activeCountdowns = new Map<string, NodeJS.Timeout>();

module.exports = {
    name: Events.VoiceStateUpdate,
    async execute(oldState, newState) {
        if (newState.channelId === null && oldState.channelId === null) return;

        const joinedChannelId = newState.channelId;
        const leftChannelId = oldState.channelId;

        const relevantChannels = [
            dconfig.channels.touch2v2,
            dconfig.channels.touch3v3,
            dconfig.channels.all3v3,
            dconfig.channels.all4v4
        ];

        // Determine which channel to act on (either joined or left)
        const chId = joinedChannelId ?? leftChannelId;
        if (!relevantChannels.includes(chId)) return;

        const channel = newState.guild.channels.cache.get(chId);
        const currentMembers = channel?.members;

        let alertsId = dconfig.channels.allAlerts;
        let isTouch = false;

        if ([dconfig.channels.touch2v2, dconfig.channels.touch3v3].includes(chId)) {
            alertsId = dconfig.channels.touchAlerts;
            isTouch = true;
        }

        // If user joined, validate them
        if (joinedChannelId) {
            const res = await Request.get(`${APIEndpoints.GET_REGISTERED_PLAYER}/${newState.member?.id}`);
            if (!res.registered) {
                CacheUtil.getChannel(newState.guild, alertsId).send({
                    content: `<@${newState.member?.id}>`,
                    embeds: [EmbedUtil.create({
                        type: "no",
                        description: `You must link your Discord account with your MC account by executing the command \`/link <code>\` at <#${dconfig.channels.register}>.`,
                    })]
                });
                newState.member?.voice.setChannel(null);
                return;
            }

            Game.refreshMemberNickname(newState.member);

            if (isTouch && !res.isTouch) {
                CacheUtil.getChannel(newState.guild, alertsId).send({
                    content: `<@${newState.member?.id}>`,
                    embeds: [EmbedUtil.create({
                        type: "no",
                        description: "You cannot queue because you were last logged in as a NON-TOUCH player (PC, PlayStation, XBox or otherwise). If you still want to queue touch-only ranked bedwars, you must first log in with a touch device.",
                    })]
                });
                newState.member?.voice.setChannel(null);
                return;
            }
        }

        const requiredSizes = {
            [dconfig.channels.touch2v2]: 4,
            [dconfig.channels.touch3v3]: 6,
            [dconfig.channels.all3v3]: 6,
            [dconfig.channels.all4v4]: 8
        };

        const teamSizes = {
            [dconfig.channels.touch2v2]: [2, 2],
            [dconfig.channels.touch3v3]: [3, 2],
            [dconfig.channels.all3v3]: [3, 2],
            [dconfig.channels.all4v4]: [4, 2]
        };

        const requiredSize = requiredSizes[chId];
        const [teamSize, numTeams] = teamSizes[chId] || [];

        if (!requiredSize || !teamSize || !numTeams || !currentMembers) return;

        const alertsChannel = CacheUtil.getChannel(newState.guild, alertsId);

        if (currentMembers.size < requiredSize) {
            if (activeCountdowns.has(chId)) {
                clearTimeout(activeCountdowns.get(chId)!);
                activeCountdowns.delete(chId);
                if (alertsChannel) {
                    alertsChannel.send({
                        content: (CacheUtil.getChannel(newState.guild, chId) as VoiceChannel).members.map(member => `<@${member.id}>`).join(" "),
                        embeds: [EmbedUtil.create({
                            type: "no",
                            description: `Game queue cancelled because <@${oldState.member?.id}> left the VC.`,
                        })]
                    });
                }
            }
            return;
        }

        if (activeCountdowns.has(chId)) return;

        if (alertsChannel) {
            alertsChannel.send({
                content: (CacheUtil.getChannel(newState.guild, chId) as VoiceChannel).members.map(member => `<@${member.id}>`).join(" "),
                embeds: [EmbedUtil.create({
                    type: "yes",
                    description: `Game is full! Starting in 5 seconds if everyone stays.`,
                })]
            });
        }

        const timeout = setTimeout(async () => {
            const finalChannel = newState.guild.channels.cache.get(chId);
            const finalMembers = finalChannel?.members;

            if (finalMembers?.size === requiredSize) {
                await initGameChannel(newState.guild, finalMembers, teamSize, numTeams);
            }

            activeCountdowns.delete(chId);
        }, 5000);

        activeCountdowns.set(chId, timeout);
    },
};

const initGameChannel = async (guild: Guild, members: Collection<string, GuildMember>, teamSize: number, teamCount: number) => {
    // INIT MC GAME
    const res = await Request.get(`${APIEndpoints.GAME_CREATE}/?teamSize=${teamSize}&teamCount=${teamCount}&custom=0`);

    // INIT DISCORD GAME
    const permissionOverwrites: OverwriteResolvable[] = [
        {
            id: guild.id,
            deny: [PermissionFlagsBits.ViewChannel]
        },
    ];

    members.forEach(member => {
        permissionOverwrites.push({
            id: member.id,
            type: OverwriteType.Member,
            allow: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.Connect]
        });
    });

    const gameName = `#${res.id.slice(0, 4).toUpperCase()} | Lobby (${teamSize}v${teamSize})`;

    const gameVC = await guild.channels.create({
        name: gameName,
        type: ChannelType.GuildVoice,
        parent: dconfig.categories.games,
        permissionOverwrites: permissionOverwrites
    });

    const gameThread = await (CacheUtil.getChannel(guild, dconfig.channels.gameChat) as TextChannel).threads.create({
        name: `Game #${res.id.slice(0, 4).toUpperCase()}`,
        autoArchiveDuration: 60,
        type: ChannelType.PrivateThread,
        invitable: false,
    });

    members.forEach(member => {
        member.voice.setChannel(gameVC);
        gameThread.members.add(member.id);
    });

    const captains = pick2Captains(members);

    const game = new Game(res.id, gameThread.id, gameVC.id, "", "", teamSize, members.map(m => m.id), captains.map(m => m.id));
    gamesDB.add(gameThread.id, game);

    await game.sendIntroductionMessage();
    await game.updateCaptainPickingMessage();
};

const pick2Captains = (members: Collection<string, GuildMember>) => {
    const membersArray = Array.from(members.values());

    const firstIndex = Math.floor(Math.random() * membersArray.length);

    let secondIndex: number;
    do {
        secondIndex = Math.floor(Math.random() * membersArray.length);
    } while (secondIndex === firstIndex);

    return [membersArray[firstIndex], membersArray[secondIndex]];
};
