import { ChannelType, Collection, Events, Guild, GuildMember, GuildTextBasedChannel, OverwriteResolvable, OverwriteType, PermissionFlagsBits, PermissionOverwrites, TextChannel, VoiceState } from "discord.js";
import { dconfig } from "../config";
import { APIEndpoints, Request } from "../api";
import { EmbedUtil } from "../core/EmbedUtil";
import { CacheUtils } from "../core/CacheUtil";
import { GameUtil } from "../core/GameUtil";

var gameCount = 0;

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

        GameUtil.refreshMemberNickname(newState.member as GuildMember);

        if (isTouch && !res.isTouch) {
            CacheUtils.getChannel(newState.guild, alertsId).send({
                content: `<@${newState.member?.id}>`,
                embeds: [EmbedUtil.create({
                    type: "no",
                    description: "You canno queue because you were last logged in as a NON-TOUCH player (PC, PlayStation, XBox or otherwise). If you still want to queue touch-only ranked bedwars, you must first log in with a touch device.",
                })]
            });
            newState.member?.voice.setChannel(null);
            return;
        }

        const chMembers = newState.channel?.members;

        switch (newState.channelId) {
            case dconfig.channels.touch2v2:
                if (chMembers?.size == 4) await initGameChannel(newState.guild, chMembers, "2v2");
                break;
            case dconfig.channels.touch3v3:
                if (chMembers?.size == 6) await initGameChannel(newState.guild, chMembers, "3v3");
                break;
            case dconfig.channels.all3v3:
                if (chMembers?.size == 6) await initGameChannel(newState.guild, chMembers, "3v3");
                break;
            case dconfig.channels.all4v4:
                if (chMembers?.size == 8) await initGameChannel(newState.guild, chMembers, "4v4");
                break;
        }
    },
};

const initGameChannel = async (guild: Guild, members: Collection<string, GuildMember>, teamSize: number, teamCount: number) => {
    // INIT MC GAME
    const res = await Request.get(`${APIEndpoints.GAME_CREATE}/?teamSize=${teamSize}&teamCount=${teamCount}&custom=0${members.map(member => `&users=${member.id}`)}`);

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

    const gameName = `Game ${res.id.slice(0, 8)} | ${teamSize}v${teamSize}`

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

    gameThread.send({
        embeds: [{
            author: {
                name: `Eliagic Ranked Bedwars | Game #${res.id}`,
                icon_url: 'https://images-ext-1.discordapp.net/external/xPUGYxZAJDXj4ScgckfwI0SvwkRQDNQDTi2gF27kRNc/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/1209943786252144690/6d272ead1117efac2cf674582fceff1f.png?format=webp&quality=lossless&width=801&height=801'
            },
            description:
                `**Game #8268 | Matchmaking**\n
                **Matchmaking Type:** Captain\n
                > 🎲 Random Captains have been chosen!\n\n
                ### Captains\n
                **Team 1 Captain:** <@${captains[0].id}>\n
                **Team 2 Captain:** <@${captains[1].id}>`,
            color: 0xFFFFFF,
            footer: {
                text: `eliagic.club | <t:${Math.floor(Date.now() / 1000)}>`
            },
        }]
    })
}

const pick2Captains = (members: Collection<string, GuildMember>) => {
    const membersArray = Array.from(members.values());

    const firstIndex = Math.floor(Math.random() * membersArray.length);

    let secondIndex;
    do {
        secondIndex = Math.floor(Math.random() * membersArray.length);
    } while (secondIndex === firstIndex);

    return [membersArray[firstIndex], membersArray[secondIndex]];
}