import {IonContent, IonPage} from '@ionic/react';
import ExploreContainer from '../components/ExploreContainer';
import React from "react";
import Header from "../components/Header";

const Type: React.FC = () => {
    return (
        <IonPage>
            <Header title="Type"/>
            <IonContent fullscreen color="light">
                <ExploreContainer name="Data type"/>
            </IonContent>
        </IonPage>
    );
};

export default Type;
