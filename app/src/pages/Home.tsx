import {IonContent, IonPage} from '@ionic/react';
import ExploreContainer from '../components/ExploreContainer';
import React from "react";
import Header from "../components/Header";

const Home: React.FC = () => {
    return (
        <IonPage>
            <Header title="ax87"/>
            <IonContent fullscreen color="light">
                <ExploreContainer name="Home page"/>
            </IonContent>
        </IonPage>
    );
};

export default Home;
