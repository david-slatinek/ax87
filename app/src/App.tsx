import {Redirect, Route} from 'react-router-dom';
import {
    IonApp,
    IonIcon,
    IonLabel,
    IonRouterOutlet,
    IonTabBar,
    IonTabButton,
    IonTabs,
    setupIonicReact
} from '@ionic/react';
import {IonReactRouter} from '@ionic/react-router';
import {code, home, list} from 'ionicons/icons';
import Home from './pages/Home';
import View from './pages/View';
import Type from "./pages/Type";

/* Core CSS required for Ionic components to work properly */
import '@ionic/react/css/core.css';

/* Basic CSS for apps built with Ionic */
import '@ionic/react/css/normalize.css';
import '@ionic/react/css/structure.css';
import '@ionic/react/css/typography.css';

/* Theme variables */
import './theme/variables.css';
import React from "react";

setupIonicReact();

const App: React.FC = () => (
    <IonApp>
        <IonReactRouter>
            <IonTabs>

                <IonRouterOutlet>
                    <Route exact path="/home">
                        <Home/>
                    </Route>

                    <Route exact path="/view">
                        <View/>
                    </Route>

                    <Route exact path="/type">
                        <Type/>
                    </Route>

                    <Route exact path="/">
                        <Redirect to="/home"/>
                    </Route>
                </IonRouterOutlet>

                <IonTabBar slot="bottom">
                    <IonTabButton tab="home" href="/home">
                        <IonIcon icon={home}/>
                        <IonLabel>Home</IonLabel>
                    </IonTabButton>

                    <IonTabButton tab="view" href="/view">
                        <IonIcon icon={list}/>
                        <IonLabel>View</IonLabel>
                    </IonTabButton>

                    <IonTabButton tab="type" href="/type">
                        <IonIcon icon={code}/>
                        <IonLabel>Type</IonLabel>
                    </IonTabButton>

                </IonTabBar>
            </IonTabs>
        </IonReactRouter>
    </IonApp>
);

export default App;
