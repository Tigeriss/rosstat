import {request} from "./request";
import {Session} from "../store/session";

export interface LoginResult {
    login: string;
    role: string;
    token: string;
}

export async function login(session: Session, login: string, password: string): Promise<LoginResult> {
    return await request(session, "login", {login, password});
}
