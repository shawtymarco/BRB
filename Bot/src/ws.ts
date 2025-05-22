import fs from 'fs';
import WebSocket from 'ws';
import { jwtSecret } from './config';

export function InitiateWebsocket() {
    const ws = new WebSocket('wss://localhost:8080/ws', {
        headers: {
            Authorization: `Bearer ${jwtSecret}`,
        },
        rejectUnauthorized: false,
        cert: fs.readFileSync('../certs/client.crt'),
        key: fs.readFileSync('../certs/client.key'),
        ca: [fs.readFileSync('../certs/ca.pem')],
    });

    ws.on('open', () => {
        console.log('Connected to Minecraft server via WebSocket');
    });

    ws.on('message', async (data) => {
        try {
            let success = true;
            const payload = JSON.parse(data.toString());

            switch (payload.type) {
                case 'kick':
                    console.log(`Kick user ${payload.userId} for reason: ${payload.reason}`);
                    ws.send(JSON.stringify({
                        status: success ? 'ok' : 'error',
                        message: success ? `User ${payload.userId} kicked.` : `Failed to kick user ${payload.userId}.`
                    }));
                    break;

                case 'role':
                    console.log(`Change role for ${payload.userId}: ${payload.action} ${payload.role}`);
                    ws.send(JSON.stringify({
                        status: success ? 'ok' : 'error',
                        message: success ? `User ${payload.userId} kicked.` : `Failed to kick user ${payload.userId}.`
                    }));
                    break;
            }
        } catch (err) {
            console.error('Invalid message from server:', err);
        }
    });

    ws.on('error', (error) => {
        console.error('WebSocket error:', error);
    });
}
