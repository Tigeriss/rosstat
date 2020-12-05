import React, {useContext} from "react";
import {
    BrowserRouter,
    Switch,
    Route,
    Redirect,
} from "react-router-dom";
import {Admin} from "./page/admin";
import {AppHeader} from "./component/appHeader";
import {Login} from "./page/login";
import {Observer} from "mobx-react";
import {Session, session} from "../store/session";
import {Logout} from "./page/logout";
import {OrdersBigPage} from "./page/ordersBig";
import {OrdersPage} from "./page/orders";
import {OrdersSmallPage} from "./page/ordersSmall";
import {OrdersPalletPage} from "./page/ordersPallet";
import {ShipmentPage} from "./page/shipment";
import {ShipmentPalletPage} from "./page/shipmentPallet";
import { PalletPrint } from "./page/palletPrint";

const SessionContext = React.createContext(session);

export function useSession(): Session {
    return useContext(SessionContext);
}

function GuestRouter() {
    return <BrowserRouter>
        <Switch>
            <Route path="/login">
                <Login/>
            </Route>
            <Route path="/">
                <Redirect to="/login"/>
            </Route>
        </Switch>
    </BrowserRouter>
}

function AuthRouter() {
    return <BrowserRouter>
        <Switch>
            <Route path="/logout">
                <Logout/>
            </Route>
            <Route path="/orders/big/:id">
                <AppHeader/>
                <OrdersBigPage />
            </Route>
            <Route path="/orders/small/:id">
                <AppHeader/>
                <OrdersSmallPage />
            </Route>
            <Route path="/orders/pallet/:id/print/:num">
                <PalletPrint />
            </Route>
            <Route path="/orders/pallet/:id">
                <AppHeader/>
                <OrdersPalletPage />
            </Route>
            <Route path="/orders">
                <AppHeader/>
                <OrdersPage />
            </Route>
            <Route path="/shipment/pallet/:id">
                <AppHeader/>
                <ShipmentPalletPage />
            </Route>
            <Route path="/shipment">
                <AppHeader/>
                <ShipmentPage />
            </Route>
            <Route path="/admin">
                <AppHeader/>
                <Admin />
            </Route>
            <Route path="/">
                <Redirect to="/orders"/>
            </Route>
        </Switch>
    </BrowserRouter>;
}

export function App() {
    return <SessionContext.Provider value={session}>
        <Observer>{() => session.currentUser == null
            ? <GuestRouter/>
            : <AuthRouter/>
        }</Observer>
    </SessionContext.Provider>;
}
