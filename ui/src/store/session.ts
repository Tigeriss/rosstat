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
import {SemanticShorthandCollection} from "semantic-ui-react/dist/commonjs/generic";
import {BreadcrumbSectionProps} from "semantic-ui-react/dist/commonjs/collections/Breadcrumb/BreadcrumbSection";

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

type PalletMatches = Array<{type: BigOrdersModel, barcode: string | null}>;

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
    curPage: "orders" | "orders-big" | "orders-small" | "orders-pallet" | "shipment" | "shipment-pallet" | "admin" | "none" = "none";

    // order related
    openedOrders: Record<string, boolean> = {};
    completedBoxes: Record<number, boolean> = {};
    currentDate: string = formatDate();
    ordersToBuild: OrdersModel[] | null = null;
    currentBigOrder: BigOrdersModel[] = [];
    currentSmallOrder: BigOrdersModel[] = [];
    currentBigPalletOrder: BigPalletModel = {pallet_num: 0, types: []};
    bigPalletOrderMatches: Record<string, PalletMatches> = {};

    // shipment related
    currentShipmentId: number | null = null;
    currentShipments: ShipmentModel[] = [];
    currentShipmentPallet: ShipmentPalletModel[] = [];
    sentPallets: Record<string, boolean> = {};

    lastError: string = "";
    lastSuccess: string = "";

    breadcrumbs: SemanticShorthandCollection<BreadcrumbSectionProps> = [];

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
        if (window.location.pathname.includes("/print/") || this.currentUser?.role === "admin") {
            return;
        }
        clearInterval(this.autoUpdateInterval as any);

        this.autoUpdateInterval = setInterval(() => {
            if (this.curPage === "orders") {
                this.fetchOrdersToBuild().catch(console.error);
            }

            if (this.curPage === "orders-big") {
                this.fetchBigOrdersToBuild().catch(console.error);
            }

            if (this.curPage === "orders-pallet") {
                this.fetchBigPallet().catch(console.error);
            }

            if (this.curPage === "shipment") {
                this.fetchShipmentReady().catch(console.error);
            }
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

        const currentBigOrder = await orders.getBigOrdersToBuild(this, this.currentOrderId);

        runInAction(() => {
            this.currentBigOrder.forEach(v => this.currentBigOrder.pop());
            currentBigOrder.forEach(v => this.currentBigOrder.push(v));
        });

    }

    async fetchSmallOrdersToBuild(): Promise<void> {
        if (this.currentOrderId == null) {
            return;
        }

        const currentSmallOrder = await orders.getSmallOrdersToBuild(this, this.currentOrderId);

        runInAction(() => {
            this.currentSmallOrder.forEach(v => this.currentSmallOrder.pop());
            currentSmallOrder.forEach(v => this.currentSmallOrder.push(v));
        });
    }

    async fetchBigPallet(): Promise<void> {
        if (this.currentOrderId == null) {
            return;
        }
        this.currentBigPalletOrder = await orders.getBigPallet(this, this.currentOrderId);

        const idp = `${this.currentOrderId}-${this.currentBigPalletOrder.pallet_num}`;
        const mts = this.bigPalletOrderMatches[idp];
        if (mts?.length > 0 && mts.length !== this.currentBigPalletOrder.types.length) {
            const types: number[] = [];
            for (const t of this.currentBigPalletOrder.types) {
                if (!types.includes(t.type)) {
                    types.push(t.type);
                }
            }

            for (const t of types) {
                const newTypes = this.currentBigPalletOrder.types.filter(v => v.type === t);
                const was = mts.filter(v => v.type.type === t).length;
                const came = newTypes.length;
                if (came > was) {
                    for (let i = newTypes.length - 1; i >= was; i--) {
                        this.bigPalletOrderMatches[idp].push({
                            barcode: null,
                            type: newTypes[i],
                        })
                    }
                }
            }
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
        return await orders.finishOrders(this, this.currentOrderId, this.preparedBoxes);
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
        const idp = `${this.currentOrderId}-${this.currentBigPalletOrder.pallet_num}`;

        if (this.bigPalletOrderMatches[idp] == null) {
            this.bigPalletOrderMatches[idp] = [];
            for (const tp of this.currentBigPalletOrder.types) {
                this.bigPalletOrderMatches[idp].push({
                    barcode: null,
                    type: tp,
                })
            }
        }
    }

    matchPalletBarcode(type: number, barcode: string): boolean {
        const idp = `${this.currentOrderId}-${this.currentBigPalletOrder.pallet_num}`;

        for (const obj of this.bigPalletOrderMatches[idp]) {
            if (obj.type.type === type && obj.barcode == null) {
                obj.barcode = barcode;
                return true;
            }
        }
        return false;
    }

    async fetchShipmentReady(): Promise<void> {
        const currentShipments = await shipment.getShipmentReady(this);

        runInAction(() => {
            this.currentShipments.forEach(v => this.currentShipments.pop());
            currentShipments.forEach(v => this.currentShipments.push(v));
        });
    }

    async fetchShipmentPallet(): Promise<void> {
        if (this.currentShipmentId == null) {
            return;
        }
        const currentShipmentPallet = await shipment.getShipmentPallet(this, this.currentShipmentId);
        runInAction(() => {
            this.currentShipmentPallet.forEach(v => this.currentShipmentPallet.pop());
            currentShipmentPallet.forEach(v => this.currentShipmentPallet.push(v));
        });
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
