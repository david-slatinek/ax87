import {IonContent, IonPage} from '@ionic/react';
import React from "react";
import Header from "../components/Header";

const Type: React.FC = () => {
    return (
        <IonPage>
            <Header title="Type"/>
            <IonContent fullscreen color="light">
            </IonContent>
        </IonPage>
    );
};

export default Type;
