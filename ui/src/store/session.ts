import {makeAutoObservable, runInAction} from "mobx";
import * as auth from "../api/auth";
import * as orders from "../api/orders";
import {BigOrdersModel, OrdersModel} from "../api/orders";

const storageKey = "user";

export class User {
    login: string;
    role: string;
    token: string;

    constructor({login, role, token}: { login: string, role: string, token: string }) {
        this.login = login;
        this.role = role;
        this.token = token;
    }
}

function formatDate(d: Date = new Date()) {
    return `${d.getDate()}.${d.getMonth()}.${d.getFullYear()} ${d.getHours()}:${d.getMinutes()}`;
}

export class Session {
    private loggedUser: User | null = null;

    get currentUser(): User | null {
        return this.loggedUser;
    }

    set currentUser(val) {
        if (val == null) {
            localStorage.removeItem(storageKey);
        } else {
            localStorage.setItem(storageKey, JSON.stringify(val));
        }
        this.loggedUser = val;
    }

    currentDate: string = formatDate();

    ordersToBuild: OrdersModel[] | null = null;

    currentBigOrder: BigOrdersModel[] = [];

    constructor() {
        makeAutoObservable(this);

        setInterval(() => {
            this.currentDate = formatDate();
        }, 1000);

        const loggedUser = localStorage.getItem(storageKey);
        if (loggedUser == null) {
            return;
        }

        try {
            const user = JSON.parse(loggedUser);
            if (typeof user.login === "string" && typeof user.role === "string" && typeof user.token === "string") {
                this.loggedUser = new User(user);
            }
        } catch (ex) {
            localStorage.removeItem(storageKey);
        }

        this.fetchOrdersToBuild().catch(console.error);
    }

    async login(login: string, password: string): Promise<boolean> {
        try {
            const res = await auth.login(this, login, password);
            if (res.token != null) {
                this.currentUser = res;
                this.fetchOrdersToBuild().catch(console.error);
                return true;
            }
        } catch (ex) {
        }
        return false;
    }

    logout() {
        this.currentUser = null;
    }

    async fetchOrdersToBuild(): Promise<void> {
        this.ordersToBuild = await orders.getOrdersToBuild(this);
    }

    async fetchBigOrdersToBuild(id: number): Promise<void> {
        this.currentBigOrder = await orders.getBigOrdersToBuild(this, id);
    }

    findOrder(id: number): OrdersModel | null {
        return (this.ordersToBuild ?? []).find(o => o.id === id) ?? null;
    }
}
