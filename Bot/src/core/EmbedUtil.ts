import { EmbedBuilder, ColorResolvable, APIEmbed } from 'discord.js';
import { dconfig } from '../config';

type EmbedType = 'yes' | 'no' | 'none';

interface EmbedOptions {
    type?: EmbedType;
    title?: string;
    description?: string;
    color?: ColorResolvable;
    timestamp?: boolean;
    footer?: string;
    thumbnail?: string;
    author?: {
        name: string;
        iconURL?: string;
    };
}

export class EmbedUtil {
    static create(options: EmbedOptions): EmbedBuilder {
        const {
            type = 'none',
            title,
            description,
            color,
            timestamp = false,
            footer,
            thumbnail,
            author
        } = options;

        const emoji = type !== 'none' ? dconfig.emojis[type] + ' ' : '';
        const embed = new EmbedBuilder();

        if (title) embed.setTitle(title);
        if (description) embed.setDescription(emoji + description);
        embed.setColor(color ?? [0x80ef80, 0xd37d77][type === "yes" ? 0 : 1]);
        if (timestamp) embed.setTimestamp();
        if (footer) embed.setFooter({ text: footer });
        if (thumbnail) embed.setThumbnail(thumbnail);
        if (author) embed.setAuthor(author);

        return embed;
    }
}