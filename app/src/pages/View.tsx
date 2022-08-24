import {IonContent, IonPage} from '@ionic/react';
import ExploreContainer from '../components/ExploreContainer';
import React from "react";
import Header from "../components/Header";

const View: React.FC = () => {
    return (
        <IonPage>
            <Header title="Last 24H"/>
            <IonContent fullscreen color="light">
                <ExploreContainer name="Tab 2 page"/>
            </IonContent>
        </IonPage>
    );
};

export default View;
