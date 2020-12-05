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

    private autoUpdateInterval = null;

    preparedBoxes: string[] = [];

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

    openedOrders: Record<string, boolean> = {};

    completedBoxes: Record<number, boolean> = {};

    currentDate: string = formatDate();

    ordersToBuild: OrdersModel[] | null = null;

    currentOrderId: number | null = null;

    currentBigOrder: BigOrdersModel[] = [];
    currentSmallOrder: BigOrdersModel[] = [];

    constructor(skip = false) {
        makeAutoObservable(this);
        if (skip) {
            return;
        }

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

        this.autoUpdate();
    }

    autoUpdate() {
        clearInterval(this.autoUpdateInterval as any);
        this.fetchOrdersToBuild().catch(console.error);

        this.autoUpdateInterval = setInterval(() => {
            this.fetchOrdersToBuild().catch(console.error);
            this.fetchBigOrdersToBuild().catch(console.error);
        }, 1000) as any;
    }

    async login(login: string, password: string): Promise<boolean> {
        try {
            const res = await auth.login(this, login, password);
            if (res.token != null) {
                this.currentUser = res;
                this.autoUpdate();
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

    async fetchBigOrdersToBuild(): Promise<void> {
        if (this.currentOrderId == null) {
            return;
        }
        this.currentBigOrder = await orders.getBigOrdersToBuild(this, this.currentOrderId);
    }

    async fetchSmallOrdersToBuild(): Promise<void> {
        if (this.currentOrderId == null) {
            return;
        }
        this.currentSmallOrder = await orders.getSmallOrdersToBuild(this, this.currentOrderId);
    }

    findOrder(id: number): OrdersModel | null {
        return (this.ordersToBuild ?? []).find(o => o.id === id) ?? null;
    }

    async finishOrders(): Promise<void> {
        if (this.currentOrderId == null) {
            return Promise.reject(new Error("orderId is null"));
        }
        return  await orders.finishOrders(this, this.currentOrderId, this.preparedBoxes);
    }
}
