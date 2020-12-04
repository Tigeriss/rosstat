import React, {useEffect} from "react";
import {useSession} from "../app";

export function Logout() {
    const session = useSession();

    useEffect(() => {
        session.currentUser = null;
    });
    
    return <div>Выходи...</div>;
}
