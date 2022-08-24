import {IonContent, IonHeader, IonPage, IonTitle, IonToolbar} from '@ionic/react';
import ExploreContainer from '../components/ExploreContainer';
import React from "react";

const Home: React.FC = () => {
    return (
        <IonPage>
            <IonHeader>
                <IonToolbar color="primary">
                    <IonTitle color="light" size="large">ax87</IonTitle>
                </IonToolbar>
            </IonHeader>
            <IonContent fullscreen color="light">
                <ExploreContainer name="Home page"/>
            </IonContent>
        </IonPage>
    );
};

export default Home;
