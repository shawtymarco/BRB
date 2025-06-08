import { ChannelType, Collection, Events, Guild, GuildMember, GuildTextBasedChannel, OverwriteResolvable, OverwriteType, PermissionFlagsBits, PermissionOverwrites, TextChannel, VoiceState } from "discord.js";
import { dconfig } from "../config";
import { APIEndpoints, Request } from "../api";
import { EmbedUtil } from "../core/EmbedUtil";
import { CacheUtils } from "../core/CacheUtil";
import { Game, gamesDB } from "../core/GameCore";

module.exports = {
    name: Events.VoiceStateUpdate,
    async execute(oldState, newState: VoiceState) {
        if (newState.channelId === null) return;
        if (![dconfig.channels.touch2v2, dconfig.channels.touch3v3, dconfig.channels.all3v3, dconfig.channels.all4v4].includes(newState.channelId)) return;

        let alertsId = dconfig.channels.allAlerts;
        let isTouch = false;

        if ([dconfig.channels.touch2v2, dconfig.channels.touch3v3].includes(newState.channelId)) {
            alertsId = dconfig.channels.touchAlerts;
            isTouch = true;
        }

        const res = await Request.get(`${APIEndpoints.GET_REGISTERED_PLAYER}/${newState.member?.id}`);
        if (!res.registered) {
            CacheUtils.getChannel(newState.guild, alertsId).send({
                content: `<@${newState.member?.id}>`,
                embeds: [EmbedUtil.create({
                    type: "no",
                    description: `You must link your Discord account with your MC account by executing the command \`/link <code>\` at <#${dconfig.channels.register}>.`,
                })]
            });
            newState.member?.voice.setChannel(null);
            return;
        }

        Game.refreshMemberNickname(newState.member as GuildMember);

        if (isTouch && !res.isTouch) {
            CacheUtils.getChannel(newState.guild, alertsId).send({
                content: `<@${newState.member?.id}>`,
                embeds: [EmbedUtil.create({
                    type: "no",
                    description: "You cannot queue because you were last logged in as a NON-TOUCH player (PC, PlayStation, XBox or otherwise). If you still want to queue touch-only ranked bedwars, you must first log in with a touch device.",
                })]
            });
            newState.member?.voice.setChannel(null);
            return;
        }

        const chMembers = newState.channel?.members;

        switch (newState.channelId) {
            case dconfig.channels.touch2v2:
                if (chMembers?.size == 4) await initGameChannel(newState.guild, chMembers, 2, 2);
                break;
            case dconfig.channels.touch3v3:
                if (chMembers?.size == 6) await initGameChannel(newState.guild, chMembers, 3, 2);
                break;
            case dconfig.channels.all3v3:
                if (chMembers?.size == 6) await initGameChannel(newState.guild, chMembers, 3, 2);
                break;
            case dconfig.channels.all4v4:
                if (chMembers?.size == 8) await initGameChannel(newState.guild, chMembers, 4, 2);
                break;
        }
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
        })
    })

    const gameName = `Game ${res.id.slice(0, 6)} | ${teamSize}v${teamSize}`

    const gameVC = await guild.channels.create({
        name: gameName,
        type: ChannelType.GuildVoice,
        parent: dconfig.categories.games,
        permissionOverwrites: permissionOverwrites
    });

    const gameThread = await (CacheUtils.getChannel(guild, dconfig.channels.gameChat) as TextChannel).threads.create({
        name: gameName,
        autoArchiveDuration: 60,
        type: ChannelType.PrivateThread,
        invitable: false,
    });

    members.forEach(member => {
        member.voice.setChannel(gameVC)
        gameThread.members.add(member.id);
    });

    const captains = pick2Captains(members);

    const game = new Game(res.id, gameThread.id, gameVC.id, teamSize, members.map(m => m.id), captains.map(m => m.id));
    gamesDB.add(gameThread.id, game)

    await game.sendIntroductionMessage()
    await game.updateCaptainPickingMessage();
}

const pick2Captains = (members: Collection<string, GuildMember>) => {
    const membersArray = Array.from(members.values());

    const firstIndex = Math.floor(Math.random() * membersArray.length);

    let secondIndex: number;
    do {
        secondIndex = Math.floor(Math.random() * membersArray.length);
    } while (secondIndex === firstIndex);

    return [membersArray[firstIndex], membersArray[secondIndex]];
}