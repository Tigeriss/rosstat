import {makeAutoObservable, runInAction} from "mobx";
import * as auth from "../api/auth";
import * as orders from "../api/orders";
import * as shipment from "../api/shipment";
import * as admin from "../api/admin";
import {
    BigOrdersModel,
    BigPalletBarcodeModel, BigPalletFinishRequestModel,
    BigPalletFinishResponseModel,
    BigPalletModel,
    OrdersModel
} from "../api/orders";
import {ShipmentModel, ShipmentPalletModel} from "../api/shipment";

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
    return `${d.getDate().toString().padStart(2, "0")}.${d.getMonth().toString().padStart(2, "0")}.${d.getFullYear().toString().padStart(4, "0")} ${d.getHours().toString().padStart(2, "0")}:${d.getMinutes().toString().padStart(2, "0")}`;
}

export class Session {
    private loggedUser: User | null = null;

    private autoUpdateInterval = null;

    preparedBoxes: string[] = [];
    users: admin.User[] = [];

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

    currentOrderId: number | null = null;

    // order related
    openedOrders: Record<string, boolean> = {};
    completedBoxes: Record<number, boolean> = {};
    currentDate: string = formatDate();
    ordersToBuild: OrdersModel[] | null = null;
    currentBigOrder: BigOrdersModel[] = [];
    currentSmallOrder: BigOrdersModel[] = [];
    currentBigPalletOrder: BigPalletModel = {pallet_num: 0, types: []};
    bigPalletOrderMatches: Array<{type: BigOrdersModel, barcode: string | null}> = [];

    // shipment related
    currentShipmentId: number | null = null;
    currentShipments: ShipmentModel[] = [];
    currentShipmentPallet: ShipmentPalletModel[] = [];
    sentPallets: Record<string, boolean> = {};

    lastError: string = "";

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
        if (window.location.pathname.includes("/print/")) {
            return;
        }
        clearInterval(this.autoUpdateInterval as any);
        if (this.currentUser?.role === "admin") {
            return;
        }

        if (this.currentUser?.role === "collector") {
            this.fetchOrdersToBuild().catch(console.error);

            this.autoUpdateInterval = setInterval(() => {
                this.fetchOrdersToBuild().catch(console.error);
                this.fetchBigOrdersToBuild().catch(console.error);
                this.fetchBigPallet().catch(console.error);
            }, 1000) as any;
        }

        if (this.currentUser?.role === "storekeeper") {
            this.fetchShipmentReady().catch(console.error);

            this.autoUpdateInterval = setInterval(() => {
                this.fetchShipmentReady().catch(console.error);
            }, 1000) as any;
        }
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

    async fetchBigPallet(): Promise<void> {
        if (this.currentOrderId == null) {
            return;
        }
        this.currentBigPalletOrder = await orders.getBigPallet(this, this.currentOrderId);

        if (this.bigPalletOrderMatches.length > 0 && this.bigPalletOrderMatches.length !== this.currentBigPalletOrder.types.length) {
            // need that?
        }
    }

    findOrder(id: number): OrdersModel | null {
        return (this.ordersToBuild ?? []).find(o => o.id === id) ?? null;
    }

    async findPallet(id: number, num: number) {
        return await orders.getPrintPallet(this, id, num);
    }

    findShipment(id: number): ShipmentModel | null {
        return (this.currentShipments ?? []).find(o => o.id === id) ?? null;
    }

    async finishOrders(): Promise<void> {
        if (this.currentOrderId == null) {
            return Promise.reject(new Error("orderId is null"));
        }
        return  await orders.finishOrders(this, this.currentOrderId, this.preparedBoxes);
    }

    async requestPalletType(barcode: string): Promise<BigPalletBarcodeModel> {
        if (this.currentOrderId == null) {
            return Promise.reject(new Error("orderId is null"));
        }
        return orders.getBigPalletBarcode(this, this.currentOrderId, barcode);
    }

    async finishBigPallet(req: BigPalletFinishRequestModel): Promise<BigPalletFinishResponseModel> {
        if (this.currentOrderId == null) {
            return Promise.reject(new Error("orderId is null"));
        }
        return orders.finishBigPallet(this, this.currentOrderId, req);
    }

    clearPalletBarcode() {
        this.bigPalletOrderMatches = [];
        for (const tp of this.currentBigPalletOrder.types) {
            this.bigPalletOrderMatches.push({
                barcode: null,
                type: tp,
            })
        }
    }

    matchPalletBarcode(type: number, barcode: string): boolean {
        for (const obj of this.bigPalletOrderMatches) {
            if (obj.type.type === type && obj.barcode == null) {
                obj.barcode = barcode;
                return true;
            }
        }
        return false;
    }

    async fetchShipmentReady(): Promise<void> {
        this.currentShipments = await shipment.getShipmentReady(this);
    }

    async fetchShipmentPallet(): Promise<void> {
        if (this.currentShipmentId == null) {
            return;
        }
        this.currentShipmentPallet = await shipment.getShipmentPallet(this, this.currentShipmentId);
    }

    async finishPalletShipment(): Promise<void> {
        if (this.currentShipmentId == null) {
            return Promise.reject(new Error("orderId is null"));
        }

        return shipment.finishPalletShipment(this, this.currentShipmentId);
    }

    async fetchUsers(): Promise<void> {
        this.users = await admin.getUsers(this);
    }

    async addUser(user: admin.User): Promise<void> {
        await admin.addUser(this, user);
        this.users.push(user);
    }


    async deleteUser(login: string): Promise<void> {
        await admin.deleteUser(this, login);
        this.users = this.users.filter(u => u.login !== login);
    }
}

export const session = new Session();
