import React, {useEffect} from "react";
import {Observer} from "mobx-react";
import {Layout} from "../component/layout";

export function OrdersBigPage() {
    useEffect(() => {

    }, []);

    return <Observer>{() =>
        <Layout>
            orders big
        </Layout>
    }</Observer>;
}
