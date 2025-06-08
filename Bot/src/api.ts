import fs from "fs";
import https from "https";
import axios from "axios";
import { jwtSecret } from "./config";

const httpsAgent = new https.Agent({
    cert: fs.readFileSync('../certs/client.crt'),
    key: fs.readFileSync('../certs/client.key'),
    ca: fs.readFileSync('../certs/ca.pem'),
    rejectUnauthorized: true
});

export enum APIEndpoints {
    CONNECT = "api/connect",
    VERIFY = "api/verify",
    GET_PLAYER = "api/players",
    GET_REGISTERED_PLAYER = "api/registered-players",
    GAME_CREATE = "api/games/create"
}

export class Request {
    static async get(endpoint: string) {
        const res = await axios.get(`https://localhost:8080/${endpoint}`, {
            headers: {
                'Authorization': `Bearer ${jwtSecret}`,
            },
            httpsAgent
        });

        return res.data;
    }

    static async post(endpoint: string, data: object) {
        const res = await axios.post(
            `https://localhost:8080/${endpoint}`,
            data,
            {
                httpsAgent,
                headers: {
                    Authorization: `Bearer ${jwtSecret}`,
                    'Content-Type': 'application/json',
                },
            }
        );

        return res.data;
    }
}