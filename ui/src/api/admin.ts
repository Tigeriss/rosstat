import {Session} from "../store/session";
import {request} from "./request";

export interface User {
    login: string;
    role: string;
    password: string;
}

export async function getUsers(session: Session): Promise<User[]> {
    return await request(session, `admin/users`, {}, "GET");
}

export async function addUser(session: Session, user: User): Promise<void> {
    return await request(session, `admin/users`, user, "POST");
}

export async function deleteUser(session: Session, login: string): Promise<User[]> {
    return await request(session, `admin/users/${login}`, {}, "DELETE");
}
