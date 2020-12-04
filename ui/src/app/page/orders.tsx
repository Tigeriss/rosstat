import {Observer} from "mobx-react";
import React, {useEffect} from "react";
import { Link } from "react-router-dom";
import {useSession} from "../app";
import {Layout} from "../component/layout";

export function OrdersPage() {
    const session = useSession();

    useEffect(() => {
        session.fetchOrders();
    }, [session]);

    return <Observer>{() =>
        <Layout>
            orders
            {JSON.stringify(session.orders)}
            <Link to="/orders/big">big</Link>
            <Link to="/orders/small">small</Link>
        </Layout>
    }</Observer>;
}
