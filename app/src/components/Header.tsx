import React from "react";
import {IonHeader, IonTitle, IonToolbar} from "@ionic/react";

interface HeaderProps {
    title: string;
}

const Header: React.FC<HeaderProps> = ({title}) => {
    return (
        <IonHeader>
            <IonToolbar color="primary">
                <IonTitle color="light" size="large">{title}</IonTitle>
            </IonToolbar>
        </IonHeader>
    );
};

export default Header;